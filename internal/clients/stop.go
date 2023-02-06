package clients

import (
	pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
)

var StopChan = make(chan struct{})

func StopAgent(req *pb.StopAgentRequest) error {
	StopChan <- struct{}{}
	return nil
}
