package model

type GlobalConfig struct {
	Colored bool `json:"colored"`
	Debug   bool `json:"debug"`
	Verbose bool `json:"verbose"`
}

func NewGlobalConfig() *GlobalConfig {
	return &GlobalConfig{}
}
