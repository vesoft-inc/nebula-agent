package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	_ "github.com/vesoft-inc/nebula-agent/internal/log"
	"github.com/vesoft-inc/nebula-go/v3/nebula"

	agent "github.com/vesoft-inc/nebula-agent/pkg/client"
	pb "github.com/vesoft-inc/nebula-agent/pkg/proto"
)

func main() {
	cfg := &agent.Config{
		Addr: &nebula.HostAddr{
			Host: "localhost",
			Port: 10000,
		},
	}

	ctx := context.TODO()
	c, err := agent.New(ctx, cfg)
	if err != nil {
		log.WithError(err).Fatalln("Create agent client error.")
	}

	uploadReq := &pb.UploadFileRequest{
		Recursively:   false,
		SourcePath:    "local source path",
		TargetBackend: &pb.Backend{Storage: &pb.Backend_Local{Local: &pb.Local{Path: "local target path"}}},
	}
	log.Infoln("Start to upload file...")
	_, e := c.UploadFile(uploadReq)
	if e != nil {
		log.WithError(e).Errorln("Upload file failed.")
	}
	log.Infoln("Upload file successfully.")
}
