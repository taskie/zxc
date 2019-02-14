package reporter

import (
	jobm "github.com/taskie/zxc/job"
	"github.com/taskie/zxc/model"
)

type FileReporter struct {
	*ReporterConfig
	*model.GlobalConfig
}

func (report *FileReporter) Report(job jobm.Job, jobResult jobm.JobResult) error {
	return nil
}
