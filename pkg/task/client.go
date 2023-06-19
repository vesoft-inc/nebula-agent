package task

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(ctx context.Context, addr string) (proto.TaskServiceClient, error) {
	logrus.Debugf("Dialing to address %s.", addr)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial agent service: %w", err)
	}
	return proto.NewTaskServiceClient(conn), nil
}
