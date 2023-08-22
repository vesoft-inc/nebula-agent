package server

import (
	"context"
	"fmt"

	"github.com/vesoft-inc/nebula-agent/v3/internal/clients"
	pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
)

// AgentServer act as an agent to interactive with services in agent machine
type AgentServer struct {
	meta *clients.NebulaMeta
}

func NewAgent(metaConfig *clients.MetaConfig) (*AgentServer, error) {
	metaclient, err := clients.NewMeta(metaConfig)
	if err != nil {
		return nil, err
	}

	a := &AgentServer{
		meta: metaclient,
	}
	return a, nil
}

// StartService start metad/storaged/graphd/all service in agent machine
func (a *AgentServer) StartService(ctx context.Context, req *pb.StartServiceRequest) (*pb.StartServiceResponse, error) {
	resp := &pb.StartServiceResponse{}

	d, err := clients.NewDaemon(clients.FromStartReq(req))
	if err != nil {
		return resp, fmt.Errorf("create service daemon failed when start service: %w", err)
	}

	return resp, d.Start()
}

// StartService stop metad/storaged/graphd/all service in agent machine
func (a *AgentServer) StopService(ctx context.Context, req *pb.StopServiceRequest) (*pb.StopServiceResponse, error) {
	resp := &pb.StopServiceResponse{}

	d, err := clients.NewDaemon(clients.FromStopReq(req))
	if err != nil {
		return resp, fmt.Errorf("create service daemon failed when stop service: %w", err)
	}

	return resp, d.Stop()
}

// ServiceStatus return the status(exit or running) of metad/storaged/graphd/all service in agent machine
func (a *AgentServer) ServiceStatus(ctx context.Context, req *pb.ServiceStatusRequest) (*pb.ServiceStatusResponse, error) {
	resp := &pb.ServiceStatusResponse{
		Status: pb.Status_UNKNOWN_STATUS,
	}

	d, err := clients.NewDaemon(clients.FromStatusReq(req))
	if err != nil {
		return resp, fmt.Errorf("create service daemon failed when get service status: %w", err)
	}

	resp.Status, err = d.Status()
	if err != nil {
		return resp, fmt.Errorf("get %s status by daemon failed: %w", req.Role, err)
	}

	return resp, nil
}

// TODO(spw): should call graphd's corresponding interface
func (a *AgentServer) BanReadWrite(context.Context, *pb.BanReadWriteRequest) (*pb.BanReadWriteResponse, error) {
	return nil, nil
}

// TODO(spw): should call graphd's corresponding interface
func (a *AgentServer) AllowReadWrite(context.Context, *pb.AllowReadWriteRequest) (*pb.AllowReadWriteResponse, error) {
	return nil, nil
}

func (a *AgentServer) DataPlayBack(ctx context.Context, req *pb.DataPlayBackRequest) (*pb.DataPlayBackResponse, error) {
	resp := &pb.DataPlayBackResponse{}

	return resp, clients.NewPlayBack(req).PlayBack()
}

func (a *AgentServer) StopAgent(ctx context.Context, req *pb.StopAgentRequest) (*pb.StopAgentResponse, error) {
	resp := &pb.StopAgentResponse{}
	return resp, clients.StopAgent(req)
}

func (a *AgentServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: "healthy",
	}, nil
}

func (a *AgentServer) GetSpaceUsages(ctx context.Context, req *pb.GetSpaceUsagesRequest) (*pb.GetSpaceUsagesResponse, error) {
	return clients.NewSpaceUsage(req.DataPath).GetSpaceUsages()
}
