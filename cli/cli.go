package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/taskie/zxc/zxcrpc"

	"github.com/taskie/zxc/job"
	"github.com/taskie/zxc/model"

	jobm "github.com/taskie/zxc/job"

	flags "github.com/jessevdk/go-flags"
	isatty "github.com/mattn/go-isatty"
	log "github.com/sirupsen/logrus"
	"github.com/taskie/zxc"
)

var (
	version = zxc.Version
)

type Options struct {
	Config        string   `short:"c" long:"config" description:"zxc config"`
	JobTemplate   string   `short:"t" long:"jobTemplate" description:"job template"`
	JobParamsJSON string   `short:"p" long:"jobParamsJson" description:"job parameters (JSON)"`
	JobParams     []string `short:"P" long:"jobParam" description:"job parameter (Key=Value)"`
	Stdin         string   `short:"i" long:"stdin" description:"redirect standard input to this file"`
	Stdout        string   `short:"o" long:"stdout" description:"standard output to this file"`
	Stderr        string   `short:"e" long:"stderr" description:"standard error to this file"`
	Color         func()   `long:"color" description:"colorize output"`
	NoColor       func()   `long:"noColor" description:"NOT colorize output"`
	colored       bool
	Daemon        bool   `short:"d" long:"daemon"`
	DaemonListen  string `short:"l" long:"daemonListen"`
	Debug         bool   `long:"debug" description:"debug mode"`
	Verbose       bool   `short:"v" long:"verbose" description:"show verbose output"`
	Version       bool   `short:"V" long:"version" description:"show version"`
}

func run(args []string, opts Options) (int, error) {
	if opts.Daemon {
		zxc := zxcrpc.NewZxcRPCServer()
		listen := opts.DaemonListen
		if listen == "" {
			listen = ":3026"
		}
		err := zxc.Serve(listen)
		if err != nil {
			return 1, err
		}
		return 0, nil
	}

	var err error
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-signalChan
		cancel()
	}()
	var config *zxc.Config
	if opts.Config != "" {
		config, err = zxc.NewConfigWithJsonFile(opts.Config)
		if err != nil {
			return 0, err
		}
	} else {
		config = &zxc.Config{}
	}
	if config.Global == nil {
		config.Global = &model.GlobalConfig{}
	}
	config.Global.Colored = opts.colored
	config.Global.Debug = opts.Debug
	config.Global.Verbose = opts.Verbose
	app := zxc.NewZxc(config)

	var jobTemplateConfig *jobm.JobTemplateConfig
	if opts.JobTemplate != "" {
		jobTemplateConfig, err = job.NewJobTemplateConfigWithJsonFile(opts.JobTemplate)
		if err != nil {
			return 127, err
		}
	} else {
		jobTemplateConfig = &job.JobTemplateConfig{}
	}
	if args != nil && len(args) > 1 {
		jobTemplateConfig.Command = args[1:]
	}
	if opts.Stdin == "-" {
		jobTemplateConfig.Stdin = ""
	} else if opts.Stdin != "" {
		jobTemplateConfig.Stdin = opts.Stdin
	}
	if opts.Stdout == "-" {
		jobTemplateConfig.Stdout = ""
	} else if opts.Stdout != "" {
		jobTemplateConfig.Stdout = opts.Stdout
	}
	if opts.Stderr == "-" {
		jobTemplateConfig.Stderr = ""
	} else if opts.Stderr != "" {
		jobTemplateConfig.Stderr = opts.Stderr
	}

	var jobConfig *jobm.JobConfig
	jobTemplate := jobm.NewJobTemplate(jobTemplateConfig)
	if opts.JobParamsJSON != "" {
		jobConfig, err = jobTemplate.BuildJobConfigWithJsonFile(opts.JobParamsJSON)
		if err != nil {
			return 127, err
		}
	} else {
		paramsMap := make(map[string]string)
		for _, v := range opts.JobParams {
			kv := strings.SplitN(v, "=", 2)
			if len(kv) != 2 {
				return 127, fmt.Errorf("invalid param: %s", v)
			}
			paramsMap[kv[0]] = kv[1]
		}
		jobConfig, err = jobTemplate.BuildJobConfig(paramsMap)
		if err != nil {
			return 127, err
		}
	}

	job := jobm.NewJob(jobConfig)
	jobResult, err := app.RunJob(ctx, job)
	exitStatus := 127
	if jobResult != nil && jobResult.ExitStatus != 0 {
		exitStatus = jobResult.ExitStatus
	}
	return exitStatus, err
}

func Main() {
	var opts Options

	outFd := os.Stdout.Fd()
	colored := isatty.IsTerminal(outFd) || isatty.IsCygwinTerminal(outFd)
	opts.Color = func() {
		colored = true
	}
	opts.NoColor = func() {
		colored = false
	}
	opts.colored = colored

	args, err := flags.ParseArgs(&opts, os.Args)
	if opts.Version {
		if opts.Verbose {
			fmt.Println("Version:    ", version)
		} else {
			fmt.Println(version)
		}
		os.Exit(0)
	}
	if err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	exitStatus, err := run(args, opts)
	if err != nil {
		log.Error(err)
		os.Exit(exitStatus)
	}
}
