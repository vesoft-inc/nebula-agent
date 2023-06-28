package main

import (
	"flag"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/vesoft-inc/nebula-agent/v3/internal/clients"
	"github.com/vesoft-inc/nebula-agent/v3/internal/limiter"
	_ "github.com/vesoft-inc/nebula-agent/v3/internal/log"
	"github.com/vesoft-inc/nebula-agent/v3/internal/server"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/config"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/plugin"
	pb "github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
)

var (
	GitInfoSHA string
)
var (
	agent      = flag.String("agent", "auto", "The agent server address")
	meta       = flag.String("meta", "", "The nebula metad service address, any metad address will be ok")
	hbs        = flag.Int("hbs", 60, "Agent heartbeat interval to nebula meta, in seconds")
	debug      = flag.Bool("debug", false, "Open debug will output more detail info")
	ratelimit  = flag.Int("ratelimit", 0, "Limit the file upload and download rate, unit Mbps")
	configFile = flag.String("f", "./etc/config.yaml", "the config file path")
)

func main() {
	flag.Parse()
	log.WithField("version", GitInfoSHA).Info("Start agent server...")
	config.InitConfig(*configFile)
	if *agent != "auto" {
		config.C.Agent = *agent
	}
	if config.C.Agent == "auto" {
		config.C.Agent = net.IPv4bcast.String() + ":8888"
	}
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// load plugin
	plugin.Load()

	// set agent rate limit
	limiter.Rate.SetLimiter(*ratelimit)

	lis, err := net.Listen("tcp", config.C.Agent)
	if err != nil {
		log.WithError(err).Fatalf("Failed to listen: %v.", *agent)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// graceful stop
	go func() {
		<-clients.StopChan
		log.Infoln("Stopping server...")
		grpcServer.GracefulStop()
		log.Infoln("Stopping plugins.")
		plugin.Stop()
	}()

	var agentServer *server.AgentServer
	if *meta != "" {
		metaCfg, err := clients.NewMetaConfig(*agent, *meta, GitInfoSHA, *hbs)
		if err != nil {
			log.WithError(err).Fatalf("Failed to create meta config.")
		}
		agentServer, err = server.NewAgent(metaCfg)
		if err != nil {
			log.WithError(err).Fatalf("Failed to create agent server.")
		}
	}
	var taskServer *server.TaskServer
	pb.RegisterTaskServiceServer(grpcServer, taskServer)
	pb.RegisterAgentServiceServer(grpcServer, agentServer)
	pb.RegisterStorageServiceServer(grpcServer, server.NewStorage())

	grpcServer.Serve(lis)
}
