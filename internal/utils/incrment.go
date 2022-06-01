package utils

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

const (
	WalExt          = ".wal"
	CommitLogIdName = "commitlog.id"
)

func LoadIncrementalFiles(srcDir string, commitLogId int64) ([]string, error) {
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}

	// splitId is the largest id which smaller than commitLogId
	var splitId int64 = -1
	// hasLogId check if the src dir have commitlog.id
	hasLogId := false
	walMap := make(map[string]int64)

	for _, entry := range entries {
		wid, err := ParseWalName(entry.Name())
		if err != nil {
			if entry.Name() == CommitLogIdName {
				hasLogId = true
			}
			continue
		}
		walMap[entry.Name()] = wid
		if wid <= commitLogId && splitId < wid {
			splitId = wid
		}
	}

	if splitId == -1 {
		return nil, fmt.Errorf("%s dir doesn't have incremental wal, commitlog.id is %d", srcDir, commitLogId)
	}

	filter := make([]string, 0)
	if !hasLogId {
		return nil, fmt.Errorf("%s dir doesn't have a commitlog.id file", srcDir)
	}
	// commitlog.id need upload
	filter = append(filter, CommitLogIdName)

	for name, id := range walMap {
		if id >= splitId {
			filter = append(filter, name)
		}
	}

	return filter, nil
}

func ParseWalName(name string) (int64, error) {
	if path.Ext(name) != WalExt {
		return -1, fmt.Errorf("%s  is not a .wal format", name)
	}
	strId := strings.Split(name, ".")[0]
	return strconv.ParseInt(strId, 10, 64)
}
