package main

import (
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

type KafkaClientConf struct {
	SubscribeTopic   string `required:"true"         yaml:"SubscribeTopic"`
	BootstrapServers string `yaml:"BootstrapServers"`
	GroupID          string `yaml:"GroupID"`
	AutoOffsetReset  string `yaml:"AutoOffsetReset"`
}

type Config struct {
	KafkaClient KafkaClientConf `required:"true" yaml:"KafkaClient"`
}

func ParseConfig(cfg any, configPath ...string) error {
	_ = godotenv.Load()

	configorLoader := configor.New(&configor.Config{
		Silent:               true,
		ErrorOnUnmatchedKeys: true,
		Environment:          "",
		ENVPrefix:            "-",
		Debug:                false,
		Verbose:              false,
		AutoReload:           false,
		AutoReloadInterval:   0,
		AutoReloadCallback:   nil,
	})

	if err := configorLoader.Load(cfg, configPath...); err != nil {
		return fmt.Errorf("loading env %w", err)
	}

	return nil
}
