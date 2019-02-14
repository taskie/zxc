package job

import (
	"bytes"
	"encoding/json"
	"os"
	"text/template"
	"time"

	"github.com/google/uuid"
)

type JobTemplateConfig struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Command          []string `json:"command"`
	WorkingDirectory string   `json:"working_directory"`
	JobResultFile    string   `json:"job_result_file"`
	Stdin            string   `json:"stdin"`
	Stdout           string   `json:"stdout"`
	Stderr           string   `json:"stderr"`
	Header           string   `json:"header"`
	Footer           string   `json:"footer"`
}

func NewJobTemplateConfigWithJsonFile(jsonFile string) (*JobTemplateConfig, error) {
	r, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	dec := json.NewDecoder(r)
	var config JobTemplateConfig
	err = dec.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type JobTemplate struct {
	*JobTemplateConfig
}

func NewJobTemplate(config *JobTemplateConfig) *JobTemplate {
	if config == nil {
		config = &JobTemplateConfig{}
	}
	if config.ID == "" {
		config.ID = "{{ uuid }}"
	}
	return &JobTemplate{
		config,
	}
}

var funcMap = template.FuncMap{
	"uuid": func() string {
		return uuid.New().String()
	},
	"now": func() time.Time {
		return time.Now()
	},
	"unixTime": func(t time.Time) int64 {
		return t.Unix()
	},
	"unixTimeNano": func(t time.Time) int64 {
		return t.UnixNano()
	},
	"formattedTime": func(t time.Time, layout string) string {
		return t.Format(layout)
	},
	"rfc3339": func(t time.Time) string {
		return t.Format(time.RFC3339)
	},
	"rfc3339Nano": func(t time.Time) string {
		return t.Format(time.RFC3339Nano)
	},
	"hourly": func(t time.Time) string {
		return t.Format("2006-01-02T15")
	},
	"daily": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
	"monthly": func(t time.Time) string {
		return t.Format("2006-01")
	},
	"yearly": func(t time.Time) string {
		return t.Format("2006")
	},
}

func (jt *JobTemplate) fillString(s string, data interface{}) (string, error) {
	t := template.New("").Funcs(funcMap)
	buf := bytes.Buffer{}
	t, err := t.Parse(s)
	if err != nil {
		return "", err
	}
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}

func (jt *JobTemplate) fillStrings(ss []string, data interface{}) ([]string, error) {
	xs := make([]string, 0, 0)
	for _, s := range ss {
		x, err := jt.fillString(s, data)
		if err != nil {
			return nil, err
		}
		xs = append(xs, x)
	}
	return xs, nil
}

func (jt *JobTemplate) BuildJobConfigWithJsonFile(jsonFile string) (*JobConfig, error) {
	r, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	dec := json.NewDecoder(r)
	var data interface{}
	err = dec.Decode(&data)
	if err != nil {
		return nil, err
	}
	return jt.BuildJobConfig(data)
}

func (jt *JobTemplate) BuildJobConfig(data interface{}) (*JobConfig, error) {
	errs := make([]error, 0, 0)
	id, err := jt.fillString(jt.ID, data)
	errs = append(errs, err)
	name, err := jt.fillString(jt.Name, data)
	errs = append(errs, err)
	command, err := jt.fillStrings(jt.Command, data)
	errs = append(errs, err)
	workingDirectory, err := jt.fillString(jt.WorkingDirectory, data)
	errs = append(errs, err)
	jobResultFile, err := jt.fillString(jt.JobResultFile, data)
	errs = append(errs, err)
	stdin, err := jt.fillString(jt.Stdin, data)
	errs = append(errs, err)
	stdout, err := jt.fillString(jt.Stdout, data)
	errs = append(errs, err)
	stderr, err := jt.fillString(jt.Stderr, data)
	errs = append(errs, err)
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return &JobConfig{
		ID:               id,
		Name:             name,
		Command:          command,
		WorkingDirectory: workingDirectory,
		JobResultFile:    jobResultFile,
		Stdin:            stdin,
		Stdout:           stdout,
		Stderr:           stderr,
	}, nil
}
