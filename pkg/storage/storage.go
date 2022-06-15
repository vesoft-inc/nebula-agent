package storage

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	pb "github.com/vesoft-inc/nebula-agent/pkg/proto"
)

const (
	SessionKey = "storage-session-id"
)

// Downloader download data from externalUri in ExternalStorage to the localPath
// If localPath not exist, create it;
// If it exist, overwrite it when it is file, copy to it when it is folder
type Downloader interface {
	Download(ctx context.Context, localPath, externalUri string, recursively bool) error
}

// Uploader upload data from localPath to the externalUri in ExternalStorage
// When recursively is false, upload file from localPath to externalUri
// Else upload files in localPath folder to externalUri folder
type Uploader interface {
	Upload(ctx context.Context, externalUri, localPath string, recursively bool) error
	IncrUpload(ctx context.Context, externalUri, localPath string, commitLogId, lastLogId int64) error
}

// Dir means we treat the storage organized as tree hierarchy
type Dir interface {
	ExistDir(ctx context.Context, uri string) bool
	EnsureDir(ctx context.Context, uri string, recursively bool) error
	// GetDir return Backend with specified uri and authentication information stored in storage
	GetDir(ctx context.Context, uri string) (*pb.Backend, error)
	// ListDir only list get dir names in given dir, not recursively
	ListDir(ctx context.Context, uri string) ([]string, error)
	// RemoveDir remove all files recursively
	RemoveDir(ctx context.Context, uri string) error
}

// ExternalStorage will keep the authentication information and other configuration for the storage
// Then the functions only need uri as the parameter
type ExternalStorage interface {
	Downloader
	Uploader
	Dir
}

func New(b *pb.Backend) (ExternalStorage, error) {
	log.WithField("uri", b.Uri()).Debugf("Create type: %s storage.", b.Type())

	switch b.Type() {
	case pb.LocalType:
		return &Local{}, nil
	case pb.S3Type:
		return NewS3(b)
	}
	return nil, fmt.Errorf("invalid uri: %s", b.Uri())
}
