package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/internal/utils"
	pb "github.com/vesoft-inc/nebula-agent/pkg/proto"
	"github.com/vesoft-inc/nebula-agent/pkg/storage"
	"github.com/vesoft-inc/nebula-go/v3/nebula"
)

// Client is an agent client to call agent service
// It just handle session id and is a wrapper now, and will add more action in the future
type Client interface {
	GetAddr() *nebula.HostAddr
	// dst file or dir should not exist
	UploadFile(req *pb.UploadFileRequest) (*pb.UploadFileResponse, error)
	DownloadFile(req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error)
	StartService(req *pb.StartServiceRequest) (*pb.StartServiceResponse, error)
	StopService(req *pb.StopServiceRequest) (*pb.StopServiceResponse, error)
	BanReadWrite(req *pb.BanReadWriteRequest) (*pb.BanReadWriteResponse, error)
	AllowReadWrite(req *pb.AllowReadWriteRequest) (*pb.AllowReadWriteResponse, error)
	MoveDir(req *pb.MoveDirRequest) (*pb.MoveDirResponse, error)
	RemoveDir(req *pb.RemoveDirRequest) (*pb.RemoveDirResponse, error)
}

func genSessionId() string {
	return uuid.New().String()
}

type client struct {
	ctx     context.Context
	addr    *nebula.HostAddr
	storage pb.StorageServiceClient
	agent   pb.AgentServiceClient
}

func New(ctx context.Context, cfg *Config) (Client, error) {
	addr := utils.StringifyAddr(cfg.Addr)
	log.Debugln("Dail to address ", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial agent service: %w", err)
	}

	if ctx.Value(storage.SessionKey) == nil {
		ctx = context.WithValue(context.Background(), storage.SessionKey, genSessionId())
	}

	c := &client{
		ctx:     ctx,
		addr:    cfg.Addr,
		storage: pb.NewStorageServiceClient(conn),
		agent:   pb.NewAgentServiceClient(conn),
	}

	return c, nil
}

func (c *client) GetAddr() *nebula.HostAddr {
	return c.addr
}

func (c *client) UploadFile(req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	if c.ctx.Value(storage.SessionKey) == nil {
		return nil, fmt.Errorf("missing session in context")
	}
	req.SessionId = fmt.Sprintf("%v", c.ctx.Value(storage.SessionKey))
	return c.storage.UploadFile(c.ctx, req)
}

func (c *client) DownloadFile(req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	if c.ctx.Value(storage.SessionKey) == nil {
		return nil, fmt.Errorf("missing session in context")
	}
	req.SessionId = fmt.Sprintf("%v", c.ctx.Value(storage.SessionKey))
	return c.storage.DownloadFile(c.ctx, req)
}

func (c *client) MoveDir(req *pb.MoveDirRequest) (resp *pb.MoveDirResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("agent, move dir failed: %w", err)
		}
	}()

	return c.storage.MoveDir(c.ctx, req)
}

func (c *client) RemoveDir(req *pb.RemoveDirRequest) (resp *pb.RemoveDirResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("agent, remove dir failed: %w", err)
		}
	}()

	return c.storage.RemoveDir(c.ctx, req)
}

func (c *client) StartService(req *pb.StartServiceRequest) (resp *pb.StartServiceResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("agent, start service failed: %w", err)
		}
	}()

	return c.agent.StartService(c.ctx, req)
}

func (c *client) StopService(req *pb.StopServiceRequest) (resp *pb.StopServiceResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("agent, stop service failed: %w", err)
		}
	}()

	return c.agent.StopService(c.ctx, req)
}

func (c *client) BanReadWrite(req *pb.BanReadWriteRequest) (resp *pb.BanReadWriteResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("agent, ban read write failed: %w", err)
		}
	}()

	return c.agent.BanReadWrite(c.ctx, req)
}

func (c *client) AllowReadWrite(req *pb.AllowReadWriteRequest) (resp *pb.AllowReadWriteResponse, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("agent, allow read write failed: %w", err)
		}
	}()

	return c.agent.AllowReadWrite(c.ctx, req)
}
