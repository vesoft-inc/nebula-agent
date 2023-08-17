package clients

import pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"

type SpaceUsage struct {
	DataPath string
}

func NewSpaceUsage(dataPath string) *SpaceUsage {
	return &SpaceUsage{
		DataPath: dataPath,
	}
}

func (s *SpaceUsage) GetSpaceUsages() (*pb.GetSpaceUsagesResponse, error) {
	return nil, nil
}
