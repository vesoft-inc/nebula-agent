package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	pb "github.com/vesoft-inc/nebula-agent/pkg/proto"
)

type S3 struct {
	backend *pb.Backend

	sess   *session.Session
	client *s3.S3
}

func NewS3(b *pb.Backend) (*S3, error) {
	if b.Type() != pb.S3Type {
		return nil, fmt.Errorf("bad format s3 uri: %s", b.Uri())
	}

	creds := credentials.NewStaticCredentials(b.GetS3().AccessKey, b.GetS3().SecretKey, "")
	forcePath := strings.ContainsAny(b.GetS3().GetEndpoint(), ":")
	region := "default"
	if b.GetS3().GetRegion() != "" {
		region = b.GetS3().GetRegion()
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(b.GetS3().GetEndpoint()),
		S3ForcePathStyle: aws.Bool(forcePath), // ip:port
		Credentials:      creds,
	}))

	log.WithField("region", region).
		WithField("endpoint", b.GetS3().GetEndpoint()).
		WithField("forcePath", forcePath).
		Debugf("Try to create s3 backend.")

	return &S3{
		backend: b,
		sess:    sess,
		client:  s3.New(sess),
	}, nil
}

func (s *S3) downloadToFile(file, key string) error {
	// Create the directories in the path
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0775); err != nil {
		return fmt.Errorf("ensure dir %s failed: %w", dir, err)
	}

	// Set up the local file
	fd, err := os.Create(file)
	if err != nil {
		if err != nil {
			return fmt.Errorf("create file %s failed: %w", file, err)
		}
	}
	defer fd.Close()

	// Download the file using the AWS SDK for Go
	req := &s3.GetObjectInput{
		Bucket: aws.String(s.backend.GetS3().Bucket),
		Key:    aws.String(key),
	}
	downloader := s3manager.NewDownloader(s.sess)
	numBytes, err := downloader.Download(fd, req)
	if err != nil {
		return fmt.Errorf("download from %s to %s failed: %w", key, file, err)
	}

	log.Debugf("Download from %s to %s successfully, bytes=%d.", key, file, numBytes)
	return nil
}

func (s *S3) downloadPrefix(localDir, prefix string) error {
	req := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.backend.GetS3().GetBucket()),
		Prefix: aws.String(prefix),
	}

	err := s.client.ListObjectsV2Pages(req, func(p *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range p.Contents {
			relPath, err := filepath.Rel(prefix, *obj.Key)
			if err != nil {
				log.WithError(err).WithField("key", *obj.Key).WithField("prefix", prefix).Error("Get relative path failed")
				return false
			}

			localFile := filepath.Join(localDir, relPath)
			err = s.downloadToFile(localFile, *obj.Key)
			if err != nil {
				log.WithError(err).WithField("from", *obj.Key).WithField("to", localFile).Error("Download file failed")
				return false
			}
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("download %s recursively failed: %w", s.backend.Uri(), err)
	}

	log.Debugf("Download from %s to %s successfully.", s.backend.Uri(), localDir)
	return nil
}

func (s *S3) Download(ctx context.Context, localPath, externalUri string, recursively bool) error {
	b := s.backend.DeepCopy()
	err := b.SetUri(externalUri)
	if err != nil {
		return fmt.Errorf("download, check and set s3 uri %s failed: %w", externalUri, err)
	}

	if recursively {
		return s.downloadPrefix(localPath, b.GetS3().Path)
	} else {
		return s.downloadToFile(localPath, b.GetS3().Path)
	}
}

func (s *S3) uploadToStorage(key, file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("open file %s failed: %w when upload", file, err)
	}
	defer fd.Close()

	uploader := s3manager.NewUploader(s.sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.backend.GetS3().Bucket),
		Key:    aws.String(key),
		Body:   fd,
	})
	if err != nil {
		return fmt.Errorf("upload from %s to %s failed: %w", file, key, err)
	}

	log.Debugf("Upload from %s to %s successfully.", file, key)
	return nil
}

type fileWalk chan string

func (f fileWalk) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("walk to %s failed: %w", path, err)
	}
	if !info.IsDir() {
		f <- path
	}
	return nil
}

func (s *S3) uploadPrefix(prefix, localDir string) error {
	walker := make(fileWalk)
	go func() {
		// Gather the files to upload by walking the path recursively
		if err := filepath.Walk(localDir, walker.Walk); err != nil {
			log.WithError(err).Error("Walk failed.")
		}
		close(walker)
	}()

	// For each file found walking, upload it to S3
	// TODO(spw): create upload for each file may be ineffective
	for path := range walker {
		rel, err := filepath.Rel(localDir, path)
		if err != nil {
			log.WithError(err).Errorf("Unable to get relative path: %s.", path)
		}
		key := filepath.Join(prefix, rel)

		err = s.uploadToStorage(key, path)
		if err != nil {
			return fmt.Errorf("upload from %s to %s failed: %w", path, key, err)
		}
	}

	log.Debugf("Upload from %s to %s recursively.", localDir, s.backend.Uri())
	return nil
}

func (s *S3) Upload(ctx context.Context, externalUri, localPath string, recursively bool) error {
	b := s.backend.DeepCopy()
	err := b.SetUri(externalUri)
	if err != nil {
		return fmt.Errorf("upload, check and set s3 uri %s failed: %w", externalUri, err)
	}

	if recursively {
		return s.uploadPrefix(b.GetS3().Path, localPath)
	} else {
		return s.uploadToStorage(b.GetS3().Path, localPath)
	}
}

func (s *S3) ExistDir(ctx context.Context, uri string) bool {
	b := s.backend.DeepCopy()
	err := b.SetUri(uri)
	if err != nil {
		log.WithError(err).WithField("uri", uri).Error("Check and set uri failed when test ExistDir.")
		return false
	}

	req := &s3.ListObjectsV2Input{
		Bucket:  aws.String(b.GetS3().GetBucket()),
		Prefix:  aws.String(b.GetS3().GetPath()),
		MaxKeys: aws.Int64(1),
	}

	exist := false
	s.client.ListObjectsV2Pages(req, func(p *s3.ListObjectsV2Output, lastPage bool) bool {
		if len(p.Contents) != 0 {
			exist = true
		}
		return false
	})

	log.WithField("uri", uri).Debugf("Test exist dir %s: %v.", uri, exist)
	return exist
}

func (s *S3) EnsureDir(ctx context.Context, uri string, recursively bool) error {
	b := s.backend.DeepCopy()
	err := b.SetUri(uri)
	if err != nil {
		return fmt.Errorf("ensure dir, check and set s3 uri %s failed: %w", uri, err)
	}

	// there is no directory in s3 object store, then need do nothing
	log.WithField("uri", uri).Debugf("Ensure %s successfully.", uri)
	return nil
}

func (s *S3) GetDir(ctx context.Context, uri string) (*pb.Backend, error) {
	b := s.backend.DeepCopy()
	err := b.SetUri(uri)
	if err != nil {
		return nil, fmt.Errorf("get dir, check and set s3 uri %s failed: %w", uri, err)
	}

	log.WithField("uri", uri).Debugf("Get backend for %s successfully.", uri)
	return b, nil
}

func (s *S3) ListDir(ctx context.Context, uri string) ([]string, error) {
	b := s.backend.DeepCopy()
	err := b.SetUri(uri)
	if err != nil {
		return nil, fmt.Errorf("list dir, check and set s3 uri %s failed: %w.", uri, err)
	}

	prefix := b.GetS3().GetPath()
	if !strings.HasSuffix(prefix, "/") {
		// without "/",  we could only get current prefix
		prefix += "/"
	}
	req := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.backend.GetS3().GetBucket()),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}

	names := make([]string, 0)
	s.client.ListObjectsV2Pages(req, func(p *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range p.CommonPrefixes {
			var name string
			if prefix == "/" { // special case: backup root is the bucket root
				name = *obj.Prefix
			} else {
				name, err = filepath.Rel(prefix, *obj.Prefix)
				if err != nil {
					log.WithError(err).WithField("key", *obj.Prefix).WithField("prefix", prefix).Error("Get relative path failed.")
					return false
				}
			}
			names = append(names, name)
		}

		return true
	})

	log.WithField("uri", uri).WithField("dirs", names).Debugf("List all dirs with prefix %s successfully.", uri)
	return names, nil
}

func (s *S3) RemoveDir(ctx context.Context, uri string) error {
	b := s.backend.DeepCopy()
	err := b.SetUri(uri)
	if err != nil {
		return fmt.Errorf("remove dir, check and set s3 uri %s failed: %w", uri, err)
	}

	prefix := b.GetS3().GetPath()
	req := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.backend.GetS3().GetBucket()),
		Prefix: aws.String(prefix),
	}

	s.client.ListObjectsV2Pages(req, func(p *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range p.Contents {
			delReq := &s3.DeleteObjectInput{
				Bucket: aws.String(s.backend.GetS3().GetBucket()),
				Key:    obj.Key,
			}

			_, err = s.client.DeleteObject(delReq)
			if err != nil {
				log.WithError(err).WithField("key", *obj.Key).Errorf("Delete object %s failed.", *obj.Key)
				return false
			}
		}

		return true
	})

	log.WithField("uri", uri).Debugf("Remove all files with prefix %s successfully.", uri)
	return nil
}
