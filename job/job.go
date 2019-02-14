package job

import (
	"encoding/json"
	"os"
	"time"
)

type JobConfig struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Command          []string `json:"command"`
	WorkingDirectory string   `json:"working_directory"`
	JobResultFile    string   `json:"job_result_file"`
	Stdin            string   `json:"stdin"`
	Stdout           string   `json:"stdout"`
	Stderr           string   `json:"stderr"`
}

func NewJobConfigWithJsonFile(jsonFile string) (*JobConfig, error) {
	r, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	dec := json.NewDecoder(r)
	var config JobConfig
	err = dec.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type Job struct {
	*JobConfig
}

type JobResult struct {
	Config     *JobConfig `json:"config"`
	ExitStatus int        `json:"exit_status"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    time.Time  `json:"end_time"`
}

func NewJob(config *JobConfig) *Job {
	if config == nil {
		config = &JobConfig{}
	}
	return &Job{
		config,
	}
}
