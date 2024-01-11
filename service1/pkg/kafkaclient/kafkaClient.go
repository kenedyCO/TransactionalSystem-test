package kafkaclient

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Config struct {
	SubscribeTopic   string
	BootstrapServers string
	GroupID          string
	AutoOffsetReset  string
}

type Client struct {
	cfg      Config
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func New(cfg Config) *Client {
	return &Client{
		cfg: cfg,
	}
}
func (c *Client) Start(context.Context) error {
	var err error
	c.producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return err
	}

	c.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": c.cfg.BootstrapServers,
		"group.id":          c.cfg.GroupID,
		"auto.offset.reset": c.cfg.AutoOffsetReset,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	err = c.consumer.SubscribeTopics([]string{c.cfg.SubscribeTopic}, nil)
	if err != nil {
		return err
	}

	log.Println("client start")

	return nil
}

func (c *Client) ShutDown(context.Context) error {
	c.producer.Close()
	if err := c.consumer.Close(); err != nil {
		return err
	}

	log.Println("client end")

	return nil
}

func (c *Client) SendMessage(msg []byte, topic string) error {
	go func() {
		for e := range c.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)

	err := c.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
	if err != nil {
		return err
	}

	// Wait for message deliveries before shutting down
	c.producer.Flush(15 * 1000)

	return nil
}

func (c *Client) ReadMessage() (*kafka.Message, error) {
	return c.consumer.ReadMessage(5 * time.Second)
}
