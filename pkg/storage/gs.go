package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero/gcsfs"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/vesoft-inc/nebula-agent/v3/internal/limiter"
	"github.com/vesoft-inc/nebula-agent/v3/internal/utils"
	pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
)

const (
	defaultGCSPageSize = 100
	defaultTimout      = time.Second * 30
)

type GS struct {
	backend *pb.Backend
	client  *storage.Client
	gcsFS   *gcsfs.Fs
}

func NewGS(b *pb.Backend) (*GS, error) {
	ctx := context.Background()
	var opts []option.ClientOption

	if b.GetGs().Credentials != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(b.GetGs().Credentials)))
	} else if credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsFile))
	} else if credentialsJSON := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON"); credentialsJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(credentialsJSON)))
	} else if oauthAccessToken := os.Getenv("GOOGLE_OAUTH_ACCESS_TOKEN"); oauthAccessToken != "" {
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: oauthAccessToken})
		opts = append(opts, option.WithTokenSource(tokenSource))
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	c := stiface.AdaptClient(client)
	gcsFS := gcsfs.NewGcsFs(context.Background(), c)

	return &GS{
		backend: b,
		client:  client,
		gcsFS:   gcsFS,
	}, nil
}

func (g *GS) Download(ctx context.Context, localPath, externalUri string, recursively bool) error {
	b := g.backend.DeepCopy()
	err := b.SetUri(externalUri)
	if err != nil {
		return fmt.Errorf("download, check and set gs uri %s failed: %w", externalUri, err)
	}

	if recursively {
		return g.downloadPrefix(localPath, externalUri, externalUri)
	} else {
		return g.downloadToFile(localPath, b.GetGs().Path)
	}
}

func (g *GS) Upload(ctx context.Context, externalUri, localPath string, recursively bool) error {
	b := g.backend.DeepCopy()
	err := b.SetUri(externalUri)
	if err != nil {
		return fmt.Errorf("upload, check and set gs uri %s failed: %w", externalUri, err)
	}

	if recursively {
		return g.uploadPrefix(b.GetGs().Path, localPath)
	} else {
		return g.uploadToStorage(b.GetGs().Path, localPath)
	}
}

func (g *GS) IncrUpload(ctx context.Context, externalUri, localPath string, commitLogId, lastLogId int64) error {
	b := g.backend.DeepCopy()
	if err := b.SetUri(externalUri); err != nil {
		return fmt.Errorf("IncrUpload, check and set gs uri %s failed: %w", externalUri, err)
	}

	// check local path
	srcInfo, err := os.Stat(localPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("local path: %s does not exist", localPath)
	}
	if err != nil {
		return fmt.Errorf("get %s status err: %w", localPath, err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("%s is a file, must specify the partition dir", localPath)
	}

	// incremental local copy
	iNames, err := utils.LoadIncrFiles(localPath, commitLogId, lastLogId)
	if err != nil {
		return err
	}

	for _, iName := range iNames {
		dst := filepath.Join(b.GetGs().Path, iName)
		src := filepath.Join(localPath, iName)
		if err = g.uploadToStorage(dst, src); err != nil {
			return err
		}
	}

	return nil
}

func (g *GS) ExistDir(ctx context.Context, uri string) bool {
	b := g.backend.DeepCopy()
	if err := b.SetUri(uri); err != nil {
		log.WithField("uri", uri).Errorf("ExistDir, check and set uri %s failed: %v", uri, err)
		return false
	}

	bucket := b.GetGs().Bucket
	path := b.GetGs().Path
	prefix := getPrefix(path)

	if prefix == "" {
		return false
	}

	query := &storage.Query{
		Prefix:    prefix,
		Delimiter: "/",
	}
	it := g.client.Bucket(bucket).Objects(context.Background(), query)
	if _, err := it.Next(); err == nil {
		return true
	}

	return false
}

func (g *GS) EnsureDir(ctx context.Context, uri string, recursively bool) error {
	b := g.backend.DeepCopy()
	if err := b.SetUri(uri); err != nil {
		return fmt.Errorf("EnsureDir, check and set uri %s failed: %v", uri, err)
	}
	// there is no directory in GS object store, then need do nothing
	log.WithField("uri", uri).Debugf("ensure %s successfully.", uri)
	return nil
}

func (g *GS) GetDir(ctx context.Context, uri string) (*pb.Backend, error) {
	b := g.backend.DeepCopy()
	if err := b.SetUri(uri); err != nil {
		return nil, fmt.Errorf("GetDir, check and set uri %s failed: %v", uri, err)
	}
	log.WithField("uri", uri).Debugf("get backend for %s successfully.", uri)
	return b, nil
}

func (g *GS) ListDir(ctx context.Context, uri string) ([]string, error) {
	b := g.backend.DeepCopy()
	if err := b.SetUri(uri); err != nil {
		return nil, fmt.Errorf("ListDir, check and set uri %s failed: %w", uri, err)
	}

	bucket := b.GetGs().Bucket
	path := b.GetGs().Path
	prefix := getPrefix(path)

	names, err := g.readdirNames(bucket, prefix)
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (g *GS) RemoveDir(ctx context.Context, uri string) error {
	pathInfo, err := g.gcsFS.Stat(uri)
	if errors.Is(err, gcsfs.ErrFileNotFound) {
		// return early if file doesn't exist
		return nil
	}
	if err != nil {
		return err
	}

	if !pathInfo.IsDir() {
		return g.gcsFS.Remove(uri)
	}

	b := g.backend.DeepCopy()
	if err := b.SetUri(uri); err != nil {
		return fmt.Errorf("RemoveDir, check and set uri %s failed: %w", uri, err)
	}

	bucket := b.GetGs().Bucket
	path := b.GetGs().Path
	prefix := getPrefix(path)

	names, err := g.readdirNames(bucket, prefix)
	if err != nil {
		return err
	}
	for _, name := range names {
		subUri := uri + "/" + name
		err = g.RemoveDir(ctx, subUri)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GS) getObjectSize(bucket, key string) (int64, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(defaultTimout))
	defer cancel()

	bkt := g.client.Bucket(bucket)
	obj := bkt.Object(key)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return 0, err
	}
	return attrs.Size, nil
}

func (g *GS) uploadToStorage(key, file string) error {
	// bucket := "bucket-name"
	// key := "path/object-name"
	// file := "local-path/file.txt"
	if limiter.Rate.IsSet() {
		srcInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		limiter.Rate.Wait(srcInfo.Size())
	}

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("open file %s failed: %v", file, err)
	}
	defer f.Close()

	o := g.client.Bucket(g.backend.GetGs().Bucket).Object(key)
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	//attrs, err := o.Attrs(ctx)
	//if err != nil {
	//	return fmt.Errorf("object.Attrs: %w", err)
	//}
	//o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})
	wc := o.NewWriter(context.Background())
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	return nil
}

func (g *GS) uploadPrefix(prefix, localDir string) error {
	walker := make(fileWalk)
	go func() {
		if err := filepath.Walk(localDir, walker.Walk); err != nil {
			log.WithError(err).Error("Walk failed.")
		}
		close(walker)
	}()

	for path := range walker {
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return fmt.Errorf("unable to get relative path: %s", path)
		}
		key := filepath.Join(prefix, relPath)
		err = g.uploadToStorage(key, path)
		if err != nil {
			return fmt.Errorf("upload from %s to %s failed: %w", path, key, err)
		}
	}

	log.Debugf("Upload from %s to %s recursively.", localDir, g.backend.Uri())
	return nil
}

func (g *GS) downloadToFile(file, key string) error {
	// bucket := "bucket-name"
	// key := "path/object-name"
	// file := "local-path/file.txt"
	if limiter.Rate.IsSet() {
		size, err := g.getObjectSize(g.backend.GetGs().Bucket, key)
		if err != nil {
			return err
		}
		limiter.Rate.Wait(size)
	}

	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0775); err != nil {
		return fmt.Errorf("ensure dir %s failed: %w", dir, err)
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	rc, err := g.client.Bucket(g.backend.GetGs().Bucket).Object(key).NewReader(context.Background())
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %w", key, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}

	return nil
}

func (g *GS) downloadPrefix(localPath, uri, baseUri string) error {
	return g.downloadToDir(localPath, uri, baseUri)
}

// localPath: /tmp
// uri: gs://nebula-2024/BACKUP_2024_01_20_01_20_56/data/2/3/...
// baseUri: gs://nebula-2024/BACKUP_2024_01_20_01_20_56/
// object key: BACKUP_2023_11_14_04_49_54/data/2/3/data/000009.sst
func (g *GS) downloadToDir(localPath, uri, baseUri string) error {
	pathInfo, err := g.gcsFS.Stat(uri)
	if errors.Is(err, gcsfs.ErrFileNotFound) {
		// return early if file doesn't exist
		return nil
	}
	if err != nil {
		return err
	}

	if !pathInfo.IsDir() {
		relPath, err := filepath.Rel(baseUri, uri)
		if err != nil {
			return fmt.Errorf("unable to get relative path: %s", relPath)
		}
		localFile := filepath.Join(localPath, relPath)
		objectKey, err := getObjectKey(uri)
		if err != nil {
			return err
		}
		return g.downloadToFile(localFile, objectKey)
	}

	b := g.backend.DeepCopy()
	if err := b.SetUri(uri); err != nil {
		return fmt.Errorf("downloadToDir, check and set uri %s failed: %w", uri, err)
	}

	bucket := b.GetGs().Bucket
	path := b.GetGs().Path
	prefix := getPrefix(path)

	names, err := g.readdirNames(bucket, prefix)
	if err != nil {
		return err
	}
	for _, name := range names {
		subUri := uri + "/" + name
		err = g.downloadToDir(localPath, subUri, baseUri)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GS) readdirNames(bucket, prefix string) ([]string, error) {
	names := make([]string, 0)
	query := &storage.Query{
		Prefix:    prefix,
		Delimiter: "/",
	}
	it := g.client.Bucket(bucket).Objects(context.Background(), query)
	for {
		object, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Bucket(%q).Objects(): %w", bucket, err)
		}
		if object.Name != "" {
			fileEntry := filepath.Base(object.Name)
			names = append(names, fileEntry)
			continue
		}
		if object.Prefix != "" {
			// It's a virtual folder! It does not have a name, but prefix - this is how GCS API
			// deals with them at the moment
			dirEntry := filepath.Base(object.Prefix)
			if err != nil {
				return nil, fmt.Errorf("Bucket(%q).Objects(): %v", bucket, err)
			}
			names = append(names, dirEntry)
		}
	}
	return names, nil
}

func getObjectKey(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return strings.TrimLeft(u.Path, "/ "), nil
}

func getPrefix(path string) string {
	prefix := ""
	if path != "" && path != "." && path != "/" {
		prefix = strings.TrimPrefix(path, "/")
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
	}
	return prefix
}
