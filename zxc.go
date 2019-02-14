package zxc

import (
	"context"
	"encoding/json"
	"os"

	"github.com/k0kubun/pp"
	log "github.com/sirupsen/logrus"
	"github.com/taskie/zxc/client"
	"github.com/taskie/zxc/executor"
	"github.com/taskie/zxc/job"
	"github.com/taskie/zxc/model"
	"github.com/taskie/zxc/reporter"
	"github.com/taskie/zxc/retrier"
)

var (
	Version  = "0.1.0-beta"
	Revision = ""
)

type Config struct {
	Global   *model.GlobalConfig                 `json:"global"`
	Client   *client.ZxcClientConfig             `json:"client"`
	Retrier  *retrier.RetrierConfig              `json:"retrier"`
	Executor *executor.ExecutorConfig            `json:"executor"`
	Reporter map[string]*reporter.ReporterConfig `json:"reporter"`
}

func NewConfigWithJsonFile(jsonFile string) (*Config, error) {
	r, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	dec := json.NewDecoder(r)
	var config Config
	err = dec.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type Zxc struct {
	ZxcClient   client.ZxcClient
	Retrier     retrier.Retrier
	Executor    executor.Executor
	ReporterMap map[string]reporter.Reporter
	config      *Config
}

func NewZxc(config *Config) *Zxc {
	reporterMap := make(map[string]reporter.Reporter)
	if config.Reporter != nil {
		for k, reporterConfig := range config.Reporter {
			reporterMap[k] = reporter.NewReporter(reporterConfig, config.Global)
		}
	}
	var zxcClient client.ZxcClient
	if config.Client != nil {
		zxcClient = client.NewZxcClient(config.Client, config.Global)
	}
	zxc := &Zxc{
		ZxcClient:   zxcClient,
		Retrier:     retrier.NewRetrier(config.Retrier, config.Global),
		Executor:    executor.NewExecutor(config.Executor, config.Global),
		ReporterMap: reporterMap,
		config:      config,
	}
	return zxc
}

func (zxc *Zxc) RunJob(ctx context.Context, job *job.Job) (*job.JobResult, error) {
	if zxc.config.Global.Debug {
		pp.Println(zxc)
	}
	if zxc.ZxcClient != nil {
		err := zxc.ZxcClient.Dial()
		if err != nil {
			return nil, err
		}
		defer zxc.ZxcClient.Close()
	}
	for {
		if zxc.ZxcClient != nil {
			err := zxc.ZxcClient.DidStartJob(ctx, job.JobConfig)
			if err != nil {
				log.Warn(err)
			}
		}
		jr, err := zxc.Executor.Execute(ctx, job)
		if zxc.ZxcClient != nil {
			zxcErr := zxc.ZxcClient.DidEndJob(ctx, jr)
			if zxcErr != nil {
				log.Warn(zxcErr)
			}
		}
		if err == nil {
			return jr, nil
		}
		ok, _ := zxc.Retrier.CheckIfRetriableNow()
		if !ok {
			return jr, err
		}
		zxc.Retrier.Retry()
	}
}
