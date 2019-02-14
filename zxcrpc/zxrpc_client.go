package zxcrpc

import (
	"google.golang.org/grpc"
)

type ZxcRPCClientContainer struct {
	Client     ZxcRPCClient
	ClientConn *grpc.ClientConn
}

func Dial(target string) (*ZxcRPCClientContainer, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := NewZxcRPCClient(conn)
	return &ZxcRPCClientContainer{
		Client:     client,
		ClientConn: conn,
	}, nil
}

func (zcc *ZxcRPCClientContainer) Close() error {
	return zcc.ClientConn.Close()
}
