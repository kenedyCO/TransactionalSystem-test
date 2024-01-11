package main

import (
	"log"

	"service1/handler"
	"service1/pkg/kafkaclient"
	"service1/pkg/repository"
	"service1/pkg/runner"
	"service1/pkg/server"
	"service1/usecase"
)

func main() {
	var cfg Config
	if err := ParseConfig(&cfg, "./service1/config/config.yml"); err != nil {
		log.Println("Config fail ", err)
	}

	mainServer := server.New(server.Config{Port: cfg.HTTP.Port})

	mainRepository := repository.New(repository.Config{ConnString: cfg.Repository.ConnString})

	mainKafkaClient := kafkaclient.New(kafkaclient.Config{
		SubscribeTopic:   cfg.KafkaClient.SubscribeTopic,
		BootstrapServers: cfg.KafkaClient.BootstrapServers,
		GroupID:          cfg.KafkaClient.GroupID,
		AutoOffsetReset:  cfg.KafkaClient.AutoOffsetReset,
	})

	mainUsecase := usecase.New(usecase.Config{}, mainRepository, mainKafkaClient)

	mainHandler := handler.New(mainUsecase)
	mainHandler.AddRoute(mainServer.Echo)

	mainRunner := runner.New(mainServer, mainKafkaClient, mainRepository)
	mainRunner.RunnerRun()
}
