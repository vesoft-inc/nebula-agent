package server

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/golang/groupcache/lru"
	log "github.com/sirupsen/logrus"
	pb "github.com/vesoft-inc/nebula-agent/pkg/proto"
	"github.com/vesoft-inc/nebula-agent/pkg/storage"
)

// StorageServer provide service:
//  1. upload/download dir/files between the agent machine and external storage
//  2. handle dir/files in agent machine
type StorageServer struct {
	s  *lru.Cache
	mu sync.Mutex
}

func NewStorage() *StorageServer {
	return &StorageServer{
		s: lru.New(32),
	}
}

// TOOD(spw): need to check same storage info
func (ss *StorageServer) getStorage(sid string, b *pb.Backend) (storage.ExternalStorage, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	v, ok := ss.s.Get(sid)
	if ok {
		return v.(storage.ExternalStorage), nil
	}

	sto, err := storage.New(b)
	if err != nil {
		return nil, fmt.Errorf("create storage from backend %s failed: %w", b.Uri(), err)
	}
	ss.s.Add(sid, sto)

	return sto, nil
}

// UploadFile upload the file or directory recursively from agent machine to external storage
func (ss *StorageServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	log.WithFields(
		log.Fields{
			"session_id": req.GetSessionId(),
			"src":        req.GetSourcePath(),
			"dst":        req.GetTargetBackend().Uri(),
			"recursive":  req.GetRecursively(),
		},
	).Debug("Upload file to external storage")

	res := &pb.UploadFileResponse{}
	sto, err := ss.getStorage(req.GetSessionId(), req.GetTargetBackend())
	if err != nil {
		return res, err
	}

	err = sto.Upload(ctx, req.GetTargetBackend().Uri(), req.GetSourcePath(), req.GetRecursively())
	if err != nil {
		return res, err
	}

	return res, nil
}

// DownloadFile download the file or directory recursively from external storage to agent machine
func (ss *StorageServer) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	log.WithFields(
		log.Fields{
			"session_id": req.GetSessionId(),
			"src":        req.GetSourceBackend().Uri(),
			"dst":        req.GetTargetPath(),
			"recursive":  req.GetRecursively(),
		},
	).Debug("Download file to local machine")

	res := &pb.DownloadFileResponse{}
	sto, err := ss.getStorage(req.GetSessionId(), req.GetSourceBackend())
	if err != nil {
		return res, err
	}

	err = sto.Download(ctx, req.GetTargetPath(), req.GetSourceBackend().Uri(), req.GetRecursively())
	if err != nil {
		return res, err
	}

	return res, nil
}

// MoveDir rename file/dir in agent machine
func (ss *StorageServer) MoveDir(ctx context.Context, req *pb.MoveDirRequest) (*pb.MoveDirResponse, error) {
	log.WithField("src", req.GetSrcPath()).WithField("dst", req.GetDstPath()).Debug("Rename dir")
	res := &pb.MoveDirResponse{}

	err := os.Rename(req.GetSrcPath(), req.GetDstPath())
	if err != nil {
		return res, err
	}

	return res, nil
}

// RemoveDir delete file/dir in agent machine
func (ss *StorageServer) RemoveDir(ctx context.Context, req *pb.RemoveDirRequest) (*pb.RemoveDirResponse, error) {
	log.WithField("path", req.GetPath()).Debug("Will delete dir")
	res := &pb.RemoveDirResponse{}

	err := os.RemoveAll(req.GetPath())
	if err != nil {
		return res, nil
	}

	return res, nil
}

// ExistDir check if file/dir in agent machine
func (ss *StorageServer) ExistDir(ctx context.Context, req *pb.ExistDirRequest) (*pb.ExistDirResponse, error) {
	log.WithField("path", req.GetPath()).Debug("Check if dir exist")
	res := &pb.ExistDirResponse{Exist: false}

	_, err := os.Stat(req.GetPath())
	if err == nil {
		res.Exist = true
		return res, nil
	}
	if os.IsNotExist(err) {
		return res, nil
	}

	return res, err
}
