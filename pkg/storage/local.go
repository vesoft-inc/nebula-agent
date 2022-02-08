package storage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	pb "github.com/vesoft-inc/nebula-agent/pkg/proto"
)

type Local struct {
}

func (l *Local) copyFile(ctx context.Context, dstPath, srcPath string) (err error) {
	// check context
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// copy
	src, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() {
		e := dst.Close()
		if err == nil {
			err = e
		}
	}()

	if _, err = io.Copy(dst, src); err != nil {
		return
	}
	return dst.Sync()
}

func createIfNotExists(dir string, mode os.FileMode) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	if err := os.MkdirAll(dir, mode); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: %w", dir, err)
	}

	return nil
}

func IsExist(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (l *Local) copyDir(ctx context.Context, dstDir, srcDir string) error {
	// check context
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := createIfNotExists(dstDir, 0755); err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, srcEntry := range entries {
		srcPath := filepath.Join(srcDir, srcEntry.Name())
		dstPath := filepath.Join(dstDir, srcEntry.Name())

		stat, ok := srcEntry.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", srcPath)
		}

		switch srcEntry.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExists(dstPath, 0755); err != nil {
				return err
			}
			if err := l.copyDir(ctx, dstPath, srcPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			return fmt.Errorf("%s is symbolic link", srcPath)
		default:
			if err := l.copyFile(ctx, dstPath, srcPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(dstPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		if err := os.Chmod(dstPath, srcEntry.Mode()); err != nil {
			return err
		}
	}
	return nil
}

func (l *Local) Upload(ctx context.Context, externalUri, localPath string, recursively bool) error {
	// check external uri
	if pb.ParseType(externalUri) != pb.LocalType {
		return fmt.Errorf("invalid local uri: %s", externalUri)
	}
	dstPath := strings.TrimPrefix(externalUri, pb.LocalPrefix)
	_, err := os.Stat(dstPath)
	if !os.IsNotExist(err) {
		log.WithField("path", externalUri).Info("Path to upload already exist")
	}

	// check local path
	srcInfo, err := os.Stat(localPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("local path: %s does not exist", localPath)
	}
	if err != nil {
		return fmt.Errorf("get %s status err: %w", localPath, err)
	}
	if srcInfo.IsDir() && !recursively {
		return fmt.Errorf("%s is directory, must upload recursively", localPath)
	}

	// local copy
	if srcInfo.IsDir() {
		err = l.copyDir(ctx, dstPath, localPath)
	} else {
		err = l.copyFile(ctx, dstPath, localPath)
	}

	if err != nil {
		err = fmt.Errorf("copy from %s to %s error %w", localPath, dstPath, err)
	}

	return err
}

func (l *Local) Download(ctx context.Context, localPath, externalUri string, recursively bool) error {
	// check external uri
	if pb.ParseType(externalUri) != pb.LocalType {
		return fmt.Errorf("invalid local uri: %s", externalUri)
	}
	srcPath := strings.TrimPrefix(externalUri, pb.LocalPrefix)
	srcInfo, err := os.Stat(srcPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("source external uri: %s does not exist", externalUri)
	}
	if srcInfo.IsDir() && !recursively {
		return fmt.Errorf("%s is directory, must download recursively", externalUri)
	}

	// check local path
	isExist, err := IsExist(localPath)
	if isExist {
		log.WithField("path", localPath).Info("Path to download already exist.")
	}
	if err != nil {
		return fmt.Errorf("get local path: %s status failed: %w", localPath, err)
	}

	// local copy
	if srcInfo.IsDir() {
		err = l.copyDir(ctx, localPath, srcPath)
	} else {
		err = l.copyFile(ctx, localPath, srcPath)
	}

	if err != nil {
		err = fmt.Errorf("copy from %s to %s error %w", srcPath, localPath, err)
	}

	return err
}

func (l *Local) ExistDir(ctx context.Context, uri string) bool {
	if pb.ParseType(uri) != pb.LocalType {
		log.Errorf("Invalid local uri type: %s.", uri)
		return false
	}

	p := strings.TrimPrefix(uri, pb.LocalPrefix)
	exist, err := IsExist(p)
	if err != nil {
		log.WithError(err).Errorf("Check exist failed: %s.", uri)
		return false
	}
	return exist
}

func (l *Local) EnsureDir(ctx context.Context, uri string, recursively bool) error {
	if pb.ParseType(uri) != pb.LocalType {
		return fmt.Errorf("invalid local uri type: %s", uri)
	}
	p := strings.TrimPrefix(uri, pb.LocalPrefix)

	exist, err := IsExist(p)
	if err != nil {
		return fmt.Errorf("check path %s exist failed: %w", p, err)
	}
	if exist {
		return nil
	}

	if !recursively {
		d := path.Dir(p)
		exist, err := IsExist(d)
		if err != nil {
			return fmt.Errorf("check path %s exist failed: %w", d, err)
		}
		if !exist {
			return fmt.Errorf("%s's parent dir not found", p)
		}

		return os.Mkdir(p, os.ModePerm)
	}

	return os.MkdirAll(p, os.ModePerm)
}

func (l *Local) GetDir(ctx context.Context, uri string) (*pb.Backend, error) {
	if pb.ParseType(uri) != pb.LocalType {
		return nil, fmt.Errorf("invalid local uri type: %s", uri)
	}

	b := &pb.Backend{}
	b.SetUri(uri)
	return b, nil
}

func (l *Local) ListDir(ctx context.Context, uri string) ([]string, error) {
	if pb.ParseType(uri) != pb.LocalType {
		return nil, fmt.Errorf("invalid local uri type: %s", uri)
	}
	p := strings.TrimPrefix(uri, pb.LocalPrefix)

	entries, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("read dir %s failed: %w", uri, err)
	}

	dirs := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}

	return dirs, nil
}

func (l *Local) RemoveDir(ctx context.Context, uri string) error {
	if pb.ParseType(uri) != pb.LocalType {
		return fmt.Errorf("invalid local uri type: %s", uri)
	}
	p := strings.TrimPrefix(uri, pb.LocalPrefix)

	exist, err := IsExist(p)
	if err != nil {
		return fmt.Errorf("check %s exist failed: %w", uri, err)
	}
	if !exist {
		return fmt.Errorf("%s not found", uri)
	}

	return os.RemoveAll(p)
}
