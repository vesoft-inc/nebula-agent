package clients

import (
	"os"
	"path/filepath"
	"strconv"

	pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
)

type SpaceUsage struct {
	DataPath string
}

func NewSpaceUsage(dataPath string) *SpaceUsage {
	return &SpaceUsage{
		DataPath: dataPath,
	}
}

func (s *SpaceUsage) GetSpaceUsages() (*pb.GetSpaceUsagesResponse, error) {
	spaceUsages := make([]*pb.GetSpaceUsagesResponse_SpaceUsageItem, 0)
	dirs, err := filepath.Glob(filepath.Join(s.DataPath, "*/"))
	if err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		id, err := strconv.Atoi(filepath.Base(dir))
		if err != nil {
			continue
		}

		var size int64
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				size += info.Size()
			}
			return err
		})
		if err != nil {
			return nil, err
		}

		spaceUsages = append(spaceUsages, &pb.GetSpaceUsagesResponse_SpaceUsageItem{
			Id:    int64(id),
			Usage: size,
		})
	}

	return &pb.GetSpaceUsagesResponse{SpaceUsages: spaceUsages}, nil
}
