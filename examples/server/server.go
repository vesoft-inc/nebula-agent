package main

import (
	"context"
	"log"
	"net"

	"github.com/vesoft-inc/nebula-agent/v3/internal/clients"
	"github.com/vesoft-inc/nebula-agent/v3/pkg/proto"
	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) SendHeartbeat(ctx context.Context, req *proto.SendHeartbeatRequest) (*proto.SendHeartbeatResponse, error) {
	log.Printf("receiveHeartBeat: %v", req)
	return &proto.SendHeartbeatResponse{}, nil
}

func main() {

	lis, err := net.Listen("tcp", "127.0.0.1:8889")
	if err != nil {
		log.Fatalf("Failed to listen: %v.", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// graceful stop
	go func() {
		<-clients.StopChan
		log.Println("Stopping server...")
		grpcServer.GracefulStop()
	}()

	var server Server
	proto.RegisterServerServiceServer(grpcServer, &server)
	grpcServer.Serve(lis)
}
