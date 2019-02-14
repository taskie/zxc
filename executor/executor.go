package executor

import (
	"context"

	"github.com/taskie/zxc/job"
	"github.com/taskie/zxc/model"
)

type ExecutorConfig struct {
	StdinBuffered  bool
	StdoutBuffered bool
	StderrBuffered bool
}

type Executor interface {
	Execute(ctx context.Context, job *job.Job) (*job.JobResult, error)
}

func NewExecutor(config *ExecutorConfig, globalConfig *model.GlobalConfig) Executor {
	if config == nil {
		config = &ExecutorConfig{}
	}
	if globalConfig == nil {
		globalConfig = model.NewGlobalConfig()
	}
	return &SimpleExecutor{
		config,
		globalConfig,
	}
}
