package main

import (
	"context"
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
	res, err := c.RunTask(ctx, &proto.RunTaskRequest{
		Type: proto.TaskType_SHELL,
		Data: shell,
	})
	log.Printf("res: %v, err: %v", res, err)
}
