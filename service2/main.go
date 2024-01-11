package main

import (
	"log"

	"service2/pkg/kafkaclient"
	"service2/pkg/kafkaserver"
	"service2/pkg/runner"
)

func main() {
	var cfg Config
	if err := ParseConfig(&cfg, "./service2/config/config.yml"); err != nil {
		log.Println("Config fail ", err)
	}

	mainKafkaClient := kafkaclient.New(kafkaclient.Config{
		SubscribeTopic:   cfg.KafkaClient.SubscribeTopic,
		BootstrapServers: cfg.KafkaClient.BootstrapServers,
		GroupID:          cfg.KafkaClient.GroupID,
		AutoOffsetReset:  cfg.KafkaClient.AutoOffsetReset,
	})

	mainKafkaServer := kafkaserver.New(kafkaserver.Config{}, mainKafkaClient)

	mainRunner := runner.New(mainKafkaServer, mainKafkaClient)
	mainRunner.RunnerRun()
}
