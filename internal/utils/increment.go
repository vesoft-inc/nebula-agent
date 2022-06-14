package utils

import (
	"fmt"
	"io/ioutil"
	"math"
	"path"
	"strconv"
	"strings"
)

const (
	WalExt            = ".wal"
	CommitLogFileName = "commitlog.id"
)

// LoadIncrFiles load incremental wal file and commitlog.id file
func LoadIncrFiles(srcDir string, commitLogId, lastLogId int64) ([]string, error) {
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}
	var (
		// splitId is the largest wal start id which smaller than commitLogId
		splitId int64 = -1

		// minWalStartId is the min wal start id in src dir
		minWalStartId int64 = math.MaxInt64
	)

	// hasCommitFile check if the src dir have commitlog.id file
	hasCommitFile := false
	walMap := make(map[string]int64)

	for _, entry := range entries {
		walStartId, err := parseWalName(entry.Name())
		if err != nil {
			if entry.Name() == CommitLogFileName {
				hasCommitFile = true
			}
			continue
		}
		walMap[entry.Name()] = walStartId
		if walStartId <= commitLogId && splitId < walStartId {
			splitId = walStartId
		}
		if walStartId < minWalStartId {
			minWalStartId = walStartId
		}
	}

	if splitId == -1 {
		return nil, fmt.Errorf("%s dir doesn't have incremental wal, commitlog.id is %d", srcDir, commitLogId)
	}

	if lastLogId < minWalStartId {
		return nil, fmt.Errorf("%s dir's wal is discontinuous, lastLogId is %d and minWalStartId is %d", srcDir, lastLogId, minWalStartId)
	}

	filter := make([]string, 0)
	if !hasCommitFile {
		return nil, fmt.Errorf("%s dir doesn't have a commitlog.id file", srcDir)
	}
	// commitlog.id need upload
	filter = append(filter, CommitLogFileName)

	for name, id := range walMap {
		if id >= splitId {
			filter = append(filter, name)
		}
	}

	return filter, nil
}

func parseWalName(name string) (int64, error) {
	if path.Ext(name) != WalExt {
		return -1, fmt.Errorf("%s  is not a .wal format", name)
	}
	strId := strings.Split(name, ".")[0]
	return strconv.ParseInt(strId, 10, 64)
}
