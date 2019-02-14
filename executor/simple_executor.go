package executor

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/k0kubun/pp"
	jobm "github.com/taskie/zxc/job"
	"github.com/taskie/zxc/jsons"
	"github.com/taskie/zxc/model"
)

type SimpleExecutor struct {
	*ExecutorConfig
	*model.GlobalConfig
}

func (executor *SimpleExecutor) Execute(ctx context.Context, job *jobm.Job) (*jobm.JobResult, error) {
	if executor.Debug {
		pp.Println(job)
	}

	if len(job.Command) == 0 {
		return nil, errors.New("invalid job: command is empty")
	}

	myCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd := exec.CommandContext(myCtx, job.Command[0], job.Command[1:]...)
	if job.WorkingDirectory != "" {
		cmd.Dir = job.WorkingDirectory
	}

	if job.Stdin != "" {
		r, err := os.Open(job.Stdin)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		cmd.Stdin = r
	} else {
		cmd.Stdin = os.Stdin
	}
	if executor.StdinBuffered {
		cmd.Stdin = bufio.NewReader(cmd.Stdin)
	}

	if job.Stdout != "" {
		w, err := os.Create(job.Stdout)
		if err != nil {
			return nil, err
		}
		defer w.Close()
		cmd.Stdout = w
	} else {
		cmd.Stdout = os.Stdout
	}
	if executor.StdoutBuffered {
		cmd.Stdout = bufio.NewWriter(cmd.Stdout)
	}

	if job.Stderr != "" {
		w, err := os.Create(job.Stderr)
		if err != nil {
			return nil, err
		}
		defer w.Close()
		cmd.Stderr = w
	} else {
		cmd.Stderr = os.Stderr
	}
	if executor.StderrBuffered {
		cmd.Stderr = bufio.NewWriter(cmd.Stderr)
	}

	var jobResult jobm.JobResult
	jobResult.Config = job.JobConfig

	jobResult.StartTime = time.Now()
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	if job.JobResultFile != "" {
		jsons.EncodeToJsonFile(job.JobResultFile, &jobResult)
	}
	err := cmd.Wait()
	jobResult.EndTime = time.Now()

	if err2, ok := err.(*exec.ExitError); ok {
		if s, ok := err2.Sys().(syscall.WaitStatus); ok {
			jobResult.ExitStatus = s.ExitStatus()
		}
	}

	if executor.Debug {
		pp.Println(jobResult)
	}
	if job.JobResultFile != "" {
		jsons.EncodeToJsonFile(job.JobResultFile, &jobResult)
	}
	return &jobResult, err
}
