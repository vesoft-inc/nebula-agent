package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/vesoft-inc/nebula-agent/v3/internal/log"
	pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
)

var (
	rootDir string = "/tmp/nebula-agent"

	localFile        string = "/tmp/nebula-agent/local.txt"
	externalFile     string = "/tmp/nebula-agent/external.txt"
	localNotExist    string = "/tmp/nebula-agent/local_not_exist.txt"
	externalNotExist string = "/tmp/nebula-agent/external_not_exist.txt"

	localDir    string = "/tmp/nebula-agent/local_dir/"
	externalDir string = "/tmp/nebula-agent/external_dir/"
	resultFile  string = "/tmp/nebula-agent/result.txt"
	resultDir   string = "/tmp/nebula-agent/result/"
)

func toExternal(localname string) string {
	return pb.LocalPrefix + localname
}

func createFile(t *testing.T, filename string) {
	t.Logf("Create file %s.", filename)
	f, err := os.Create(filename)
	defer func() {
		if err = f.Close(); err != nil {
			t.Errorf("Close file %s failed: %s.", filename, err.Error())
		}
	}()

	if err != nil {
		t.Errorf("Create %s error: %s.", filename, err.Error())
	} else {
		if _, err = f.Write([]byte(filename)); err != nil {
			t.Errorf("Write content error: %s.", err.Error())
		}
	}
}

func createDir(t *testing.T, dirname string) {
	t.Logf("Create dir %s.", dirname)
	if err := os.MkdirAll(dirname, os.FileMode(0755)); err != nil {
		t.Errorf("Create dir %s failed: %s.", dirname, err.Error())
	}
}

func removeFile(t *testing.T, filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Logf("%s does not exist.", filename)
		return
	}

	if err := os.Remove(filename); err != nil {
		t.Errorf("Remove %s failed: %s.", filename, err.Error())
	}
	t.Logf("Has removed file %s.", filename)
}

func removeDir(t *testing.T, dirname string) {
	if err := os.RemoveAll(dirname); err != nil {
		t.Errorf("Remove dir %s failed: %s.", dirname, err.Error())
	}
	t.Logf("Has removed dir %s.", dirname)
}

func setup(t *testing.T) {
	// bad cases
	createDir(t, rootDir)
	createFile(t, localFile)
	createFile(t, externalFile)

	// good cases
	createDir(t, localDir)
	createFile(t, filepath.Join(localDir, "a.txt"))
	createDir(t, filepath.Join(localDir, "inner"))
	createFile(t, filepath.Join(localDir, "inner/a.txt"))

	createDir(t, externalDir)
	createFile(t, filepath.Join(externalDir, "a.txt"))
	createDir(t, filepath.Join(externalDir, "inner"))
	createFile(t, filepath.Join(externalDir, "inner/a.txt"))
}

func teardown(t *testing.T) {
	removeDir(t, rootDir)
}

func TestLocalDownload(t *testing.T) {
	type TestCase struct {
		Name        string
		ExternalUri string
		LocalPath   string
		Recursively bool
		HasError    bool
		ErrorMsg    string
	}

	testcases := []TestCase{
		{
			Name:        "source not exist",
			ExternalUri: toExternal(externalNotExist),
			LocalPath:   resultFile,
			Recursively: false,
			HasError:    true,
		},
		{
			Name:        "target exist",
			ExternalUri: toExternal(externalFile),
			LocalPath:   localFile,
			Recursively: false,
			HasError:    false,
		},
		{
			Name:        "source is dir but not recursive",
			ExternalUri: toExternal(externalDir),
			LocalPath:   resultDir,
			Recursively: false,
			HasError:    true,
		},
		{
			Name:        "normal file",
			ExternalUri: toExternal(externalFile),
			LocalPath:   resultFile,
			Recursively: false,
			HasError:    false,
		},
		{
			Name:        "normal dir",
			ExternalUri: toExternal(externalDir),
			LocalPath:   resultDir,
			Recursively: true,
			HasError:    false,
		},
	}

	setup(t)
	assert := assert.New(t)
	backend := &pb.Backend{
		Storage: &pb.Backend_Local{
			Local: &pb.Local{Path: toExternal(rootDir)},
		},
	}
	storage, err := New(backend)
	assert.Nil(err, "Create storage failed", err)

	for _, tc := range testcases {
		fmt.Printf("Test case: local download - %s\n", tc.Name)
		err = storage.Download(context.Background(), tc.LocalPath, tc.ExternalUri, tc.Recursively)
		hasError := err != nil
		assert.Equalf(hasError, tc.HasError, "Test case: '%s' failed, error: %v", tc.Name, err)
	}
	teardown(t)
}

func TestLocalUpload(t *testing.T) {
	type TestCase struct {
		Name        string
		ExternalUri string
		LocalPath   string
		Recursively bool
		HasError    bool
		ErrorMsg    string
	}

	testcases := []TestCase{
		{
			Name:        "source not exist",
			ExternalUri: toExternal(resultFile),
			LocalPath:   localNotExist,
			Recursively: false,
			HasError:    true,
		},
		{
			Name:        "target exist",
			ExternalUri: toExternal(externalFile),
			LocalPath:   localFile,
			Recursively: false,
			HasError:    false,
		},
		{
			Name:        "source is dir but not recursive",
			ExternalUri: toExternal(resultDir),
			LocalPath:   localDir,
			Recursively: false,
			HasError:    true,
		},
		{
			Name:        "normal file",
			ExternalUri: toExternal(resultFile),
			LocalPath:   localFile,
			Recursively: false,
			HasError:    false,
		},
		{
			Name:        "normal dir",
			ExternalUri: toExternal(resultDir),
			LocalPath:   localDir,
			Recursively: true,
			HasError:    false,
		},
	}

	setup(t)
	assert := assert.New(t)
	backend := &pb.Backend{
		Storage: &pb.Backend_Local{
			Local: &pb.Local{Path: pb.LocalPrefix},
		},
	}
	storage, err := New(backend)
	assert.Nilf(err, "New storage from backend failed: %v", err)
	for _, tc := range testcases {
		fmt.Printf("Test case: local upload - %s\n", tc.Name)
		err := storage.Upload(context.Background(), tc.ExternalUri, tc.LocalPath, tc.Recursively)
		hasError := err != nil
		assert.Equalf(hasError, tc.HasError, "Test case: '%s' failed, error: %v", tc.Name, err)
	}
	teardown(t)
}
