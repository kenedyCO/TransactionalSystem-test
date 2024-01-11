package kafkaserver

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Config struct {
}

type Server struct {
	cfg    Config
	client IClient
	flag   bool
}

type IClient interface {
	ReadMessage() (*kafka.Message, error)
}

func New(cfg Config, client IClient) *Server {
	return &Server{
		cfg:    cfg,
		client: client,
	}
}

func (s *Server) Start(context.Context) error {
	go func() {
		s.flag = true
		for s.flag {
			msg, err := s.client.ReadMessage()
			if err == nil {

				requestBody := msg.Value
				// Post request
				url := "http://localhost:8080/v1/service1/transaction"
				//request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
				var client http.Client
				resp, err := client.Post(url, "application/json", bytes.NewReader(requestBody))
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusOK {
					bodyBytes, err := io.ReadAll(resp.Body)
					if err != nil {
						log.Fatal(err)
					}
					bodyString := string(bodyBytes)
					log.Println(bodyString)
				}

			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				log.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}()
	log.Println("kafkaServer start")

	return nil
}

func (s *Server) ShutDown(context.Context) error {
	s.flag = false
	log.Println("kafkaServer end")

	return nil
}
