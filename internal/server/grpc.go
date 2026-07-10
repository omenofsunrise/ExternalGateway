package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "external-gateway/api/generated/gateway"
)

type GRPCServer struct {
	server  *grpc.Server
	port    int
	handler pb.ExternalGatewayServer
}

func NewGRPCServer(port int, handler pb.ExternalGatewayServer) *GRPCServer {
	grpcServer := grpc.NewServer()
	pb.RegisterExternalGatewayServer(grpcServer, handler)

	reflection.Register(grpcServer)

	return &GRPCServer{
		server:  grpcServer,
		port:    port,
		handler: handler,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Printf("gRPC server listening on port %d", s.port)
	return s.server.Serve(lis)
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
	log.Println("gRPC server stopped")
}
