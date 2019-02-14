package retrier

import "github.com/taskie/zxc/model"

type RetrierConfig struct {
	Times *int `json:"times"`
}

type Retrier interface {
	CheckIfRetriable() (bool, error)
	CheckIfRetriableNow() (bool, error)
	Retry() error
}

func NewRetrier(config *RetrierConfig, globalConfig *model.GlobalConfig) Retrier {
	times := 0
	if config != nil {
		if config.Times != nil {
			times = *config.Times
		}
	}
	return &SimpleRetrier{
		state: &SimpleRetrierState{
			LeftTimes: times,
		},
		config:       config,
		globalConfig: globalConfig,
	}
}
