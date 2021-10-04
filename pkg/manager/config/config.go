package config

import (
	"time"

	"github.com/flyteorg/flytestdlib/config"
)

//go:generate pflags Config --default-var=DefaultConfig

var (
	DefaultConfig = &Config{
		Namespace: "flyte",
		ScanInterval: config.Duration{
			Duration: 10 * time.Second,
		},
	}

	configSection = config.MustRegisterSection("manager", DefaultConfig)
)

type Config struct {
	Namespace    string  `json:"namespace" pflag:"Namespace to use for managing flytepropeller pod instances"`
	ScanInterval config.Duration `json:"scan-interval" pflag:"Frequency to scan flytepropeller pods and start / restart if necessary"`
}

func GetConfig() *Config {
	return configSection.GetConfig().(*Config)
}
