package zxcrpc

import "github.com/taskie/zxc/job"

func NewZxcRPCJob(job *job.JobConfig) *Job {
	return &Job{
		Name:    "",
		Command: job.Command,
	}
}

func NewZxcRPCJobResultMessage(jobResult *job.JobResult) *JobResultMessage {
	return &JobResultMessage{
		Client:     &Client{},
		Job:        NewZxcRPCJob(jobResult.Config),
		ExitStatus: int32(jobResult.ExitStatus),
	}
}
