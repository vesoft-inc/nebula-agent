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

	cx, cancel := context.WithDeadline(ctx, time.Now().Add(defaultTimout))
	defer cancel()

	bkt := g.client.Bucket(b.GetGs().Bucket)
	query := &storage.Query{
		Prefix:    b.GetGs().Path,
		Delimiter: "/",
	}
	it := bkt.Objects(cx, query)
	pager := iterator.NewPager(it, defaultGCSPageSize, "")

	exist := false
	for {
		var objects []*storage.ObjectAttrs
		pageToken, err := pager.NextPage(&objects)
		if err != nil {
			log.WithField("uri", uri).Errorf("ExistDir, get nex page failed: %v", err)
			return false
		}

		for _, attrs := range objects {
			if !attrs.Deleted.IsZero() {
				continue
			}
			isDir := strings.HasSuffix(attrs.Name, "/")
			if isDir && attrs.Size == 0 {
				exist = true
				break
			}
		}

		objects = nil
		if pageToken == "" {
			break
		}
	}
	log.WithField("uri", uri).Debugf("ExistDir %s: %v.", uri, exist)
	return true
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

	prefix := b.GetGs().Path
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	cx, cancel := context.WithDeadline(ctx, time.Now().Add(defaultTimout))
	defer cancel()

	names := make([]string, 0)
	bkt := g.client.Bucket(b.GetGs().Bucket)
	query := &storage.Query{
		Prefix:    prefix,
		Delimiter: "/",
	}

	it := bkt.Objects(cx, query)
	pager := iterator.NewPager(it, defaultGCSPageSize, "")

	for {
		var objects []*storage.ObjectAttrs
		pageToken, err := pager.NextPage(&objects)
		if err != nil {
			return nil, fmt.Errorf("ListDir, get nex page failed: %v", err)
		}

		for _, attrs := range objects {
			name, err := filepath.Rel(prefix, attrs.Prefix)
			if err != nil {
				return nil, fmt.Errorf("ListDir, get relative path failed: %v", err)
			}
			if strings.Contains(name, ".") {
				continue
			}
			names = append(names, name)
		}

		objects = nil
		if pageToken == "" {
			break
		}
	}
	return names, nil
}

func (g *GS) RemoveDir(ctx context.Context, uri string) error {
	if err := g.gcsFS.RemoveAll(uri); err != nil {
		if errors.Is(err, gcsfs.ErrFileNotFound) {
			return nil
		}
		return err
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
		rel, err := filepath.Rel(localDir, path)
		if err != nil {
			log.WithError(err).Errorf("Unable to get relative path: %s.", path)
		}
		key := filepath.Join(prefix, rel)

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
			return err
		}
		localFile := filepath.Join(localPath, relPath)
		objectKey, err := getObjectKey(uri)
		if err != nil {
			return err
		}
		return g.downloadToFile(localFile, objectKey)
	}

	var dir *gcsfs.GcsFile
	dir, err = g.gcsFS.Open(uri)
	if err != nil {
		return err
	}

	var infos []os.FileInfo
	infos, err = dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, info := range infos {
		subUri := uri + "/" + info.Name()
		err = g.downloadToDir(localPath, subUri, baseUri)
		if err != nil {
			return err
		}
	}
	return nil
}

func getObjectKey(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return strings.TrimLeft(u.Path, "/ "), nil
}
