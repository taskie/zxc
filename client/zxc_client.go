package client

import (
	"context"

	"github.com/taskie/zxc/job"
	"github.com/taskie/zxc/model"
)

type ZxcClientConfig struct {
	Target string
}

type ZxcClient interface {
	Dial() error
	DidStartJob(ctx context.Context, job *job.JobConfig) error
	DidEndJob(ctx context.Context, jobResult *job.JobResult) error
	Close() error
}

func NewZxcClient(config *ZxcClientConfig, globalConfig *model.GlobalConfig) ZxcClient {
	if config == nil {
		config = &ZxcClientConfig{}
	}
	if config.Target == "" {
		config.Target = "localhost:3026"
	}
	if globalConfig == nil {
		globalConfig = model.NewGlobalConfig()
	}
	return &SimpleZxcClient{
		ZxcClientConfig:       config,
		GlobalConfig:          globalConfig,
		zxcrpcClientContainer: nil,
	}
}
