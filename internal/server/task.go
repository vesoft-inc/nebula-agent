package server

import (
	"context"
	"os/exec"

	pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/task"
)

// AgentServer act as an agent to interactive with services in agent machine
type TaskServer struct{}

// StartService start metad/storaged/graphd/all service in agent machine
func (a *TaskServer) RunTask(ctx context.Context, req *pb.RunTaskRequest) (*pb.RunTaskResponse, error) {
	if req.Type == pb.TaskType_SHELL {
		return a.runShell(ctx, req)
	}
	return &pb.RunTaskResponse{
		Status: pb.TaskStatus_TASK_FAILED,
		Data:   "failed",
	}, nil
}

func (a *TaskServer) runShell(ctx context.Context, req *pb.RunTaskRequest) (*pb.RunTaskResponse, error) {
	cmd := exec.Command("bash", "-c", req.Data)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &pb.RunTaskResponse{
			Status: pb.TaskStatus_TASK_FAILED,
			Data:   err.Error(),
		}, nil
	}
	return &pb.RunTaskResponse{
		Status: pb.TaskStatus_TASK_SUCCESS,
		Data:   string(output),
	}, nil
}

func (a *TaskServer) RunStreamTask(in *pb.RunTaskRequest, out pb.TaskService_RunStreamTaskServer) error {
	return task.RunStreamShell(in.Id, in.Data, func(s string) error {
		return out.Send(&pb.StreamRunTaskResponse{
			Status: pb.TaskStatus_TASK_SUCCESS,
			Data:   s,
		})
	})
}
