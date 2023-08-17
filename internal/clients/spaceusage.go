package clients

import (
	pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
	"os"
	"path/filepath"
	"strconv"
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
	err := filepath.Walk(s.DataPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		id, err := strconv.Atoi(info.Name())
		if err != nil {
			return nil
		}

		spaceUsages = append(spaceUsages, &pb.GetSpaceUsagesResponse_SpaceUsageItem{
			Id:    int64(id),
			Usage: info.Size(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetSpaceUsagesResponse{SpaceUsages: spaceUsages}, nil
}
