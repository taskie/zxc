package client

import (
	"context"

	"github.com/taskie/zxc/job"
	"github.com/taskie/zxc/model"
	"github.com/taskie/zxc/zxcrpc"
)

type SimpleZxcClient struct {
	*ZxcClientConfig
	*model.GlobalConfig
	zxcrpcClientContainer *zxcrpc.ZxcRPCClientContainer
}

func (zxcClient *SimpleZxcClient) Dial() error {
	zcc, err := zxcrpc.Dial(zxcClient.Target)
	if err != nil {
		return err
	}
	zxcClient.zxcrpcClientContainer = zcc
	return nil
}

func (zxcClient *SimpleZxcClient) DidStartJob(ctx context.Context, job *job.JobConfig) error {
	_, err := zxcClient.zxcrpcClientContainer.Client.DidStartJob(ctx, &zxcrpc.JobMessage{
		Client: &zxcrpc.Client{},
		Job:    zxcrpc.NewZxcRPCJob(job),
	})
	return err
}

func (zxcClient *SimpleZxcClient) DidEndJob(ctx context.Context, jobResult *job.JobResult) error {
	_, err := zxcClient.zxcrpcClientContainer.Client.DidEndJob(ctx, zxcrpc.NewZxcRPCJobResultMessage(jobResult))
	return err
}

func (zxcClient *SimpleZxcClient) Close() error {
	if zxcClient.zxcrpcClientContainer != nil {
		return zxcClient.zxcrpcClientContainer.Close()
	}
	return nil
}
