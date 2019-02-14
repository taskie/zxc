package zxcrpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

type SimpleZxcRPCServer struct{}

func (server *SimpleZxcRPCServer) DidStartJob(ctx context.Context, job *JobMessage) (*Server, error) {
	log.Info("DidStartJob", job)
	return &Server{
		Name: "foo",
	}, nil
}

func (server *SimpleZxcRPCServer) DidEndJob(ctx context.Context, jobResult *JobResultMessage) (*Server, error) {
	log.Info("DidEndJob", jobResult)
	return &Server{
		Name: "foo",
	}, nil
}

func NewZxcRPCServer() *SimpleZxcRPCServer {
	return &SimpleZxcRPCServer{}
}

func (server *SimpleZxcRPCServer) Serve(address string) error {
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	RegisterZxcRPCServer(grpcServer, NewZxcRPCServer())
	return grpcServer.Serve(lis)
}
