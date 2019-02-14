package reporter

import (
	"github.com/taskie/zxc/job"
	"github.com/taskie/zxc/model"
)

type ReporterConfig struct {
}

type Reporter interface {
	Report(job job.Job, jobResult job.JobResult) error
}

func NewReporter(config *ReporterConfig, globalConfig *model.GlobalConfig) Reporter {
	return &FileReporter{
		config,
		globalConfig,
	}
}
