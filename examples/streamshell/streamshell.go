package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/task"
)

func main() {
	if len(os.Args) < 3 {
		log.Printf("usage: %s addr shell", os.Args[0])
		os.Exit(1)
	}
	addr := os.Args[1]
	shell := os.Args[2]
	log.Printf("addr: %s, shell: %s", addr, shell)
	ctx := context.TODO()
	c, err := task.New(ctx, addr)
	if err != nil {
		log.Fatalf("Create agent client error: %v", err)
	}

	// send stream shell request
	res, err := c.RunStreamTask(ctx, &proto.RunTaskRequest{
		Type: proto.TaskType_SHELL,
		Data: shell,
	})
	if err != nil {
		log.Fatalf("RunStreamTask error: %v", err)
		return
	}
	// recieve stream data
	for {
		reply, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
		}
		log.Printf("status: %v, data: %s", reply.Status, reply.Data)
	}
}
