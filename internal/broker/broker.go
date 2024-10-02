package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

const (
	brokerList = "localhost:9092"
	maxRetry   = 5
)

type Handler interface {
	Handle(requestKey []byte, requestMsg []byte) (msg OrchestratorMessage, err error)
}

type Connector interface {
	RegisterHandler(Handler)
	Run() error
	Stop() error
}

type ConnectorKafka struct {
	topicWithRequests string
	consumerGroupName string
	producer          sarama.SyncProducer
	consumer          sarama.ConsumerGroup
	handler           Handler
	stopConsumer      context.CancelFunc
	ConsumerReady     chan bool
	wg                *sync.WaitGroup
}

func (h *ConnectorKafka) Setup(_ sarama.ConsumerGroupSession) error {
	// Perform any necessary setup tasks here.
	return nil
}

func (h *ConnectorKafka) Cleanup(_ sarama.ConsumerGroupSession) error {
	// Perform any necessary cleanup tasks here.
	return nil
}

func (c *ConnectorKafka) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Got message - Topic: %s, Key: %s, Value:%s", c.topicWithRequests, string(msg.Key), string(msg.Value))
			respMsg, err := c.handler.Handle(msg.Key, msg.Value)
			if err != nil {
				return err
			}
			log.Printf("Sending message - Topic:%s, Key: %s,Value: %s", respMsg.NextTopic, respMsg.Key, respMsg.Value)
			if respMsg.NextTopic != "" {
				response := &sarama.ProducerMessage{
					Topic: respMsg.NextTopic,
					Value: sarama.StringEncoder(respMsg.Value),
					Key:   sarama.StringEncoder(respMsg.Key),
				}
				part, _, err := c.producer.SendMessage(response)
				if err != nil {
					return err
				}
				log.Printf("Sent message. Partition: %v, Topic: %s", part, respMsg.NextTopic)
			}
			session.MarkMessage(msg, "")
			return nil
		case <-session.Context().Done():
			return nil
		}
	}
}

func NewBrokerKafka(topicWithRequests string, consumerGroupName string) (Connector, error) {
	conn := &ConnectorKafka{
		topicWithRequests: topicWithRequests,
		consumerGroupName: consumerGroupName,
		ConsumerReady:     make(chan bool),
		wg:                &sync.WaitGroup{},
	}

	producer, err := configProducer()
	if err != nil {
		return nil, err
	}
	conn.producer = producer

	consumer, err := configConsumer(conn.consumerGroupName)
	if err != nil {
		return nil, err
	}
	conn.consumer = consumer

	return conn, err
}

func configProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = maxRetry
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{brokerList}, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func configConsumer(consumerGroupName string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerGroup, err := sarama.NewConsumerGroup([]string{brokerList}, consumerGroupName, config)
	if err != nil {
		return nil, err
	}
	return consumerGroup, nil
}

func (c *ConnectorKafka) RegisterHandler(h Handler) {
	c.handler = h
}

func (c *ConnectorKafka) Run() error {
	if c.handler == nil {
		return fmt.Errorf("handler is not specified")
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.stopConsumer = cancel
	c.wg.Add(1)
	log.Println("Starting consumer...")
	defer c.wg.Done()
	for {
		err := c.consumer.Consume(ctx, []string{c.topicWithRequests}, c)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("Error consuming messages: %v", err)
			return err
		}
		//c.ConsumerReady = make(chan bool)
	}

}

func (c *ConnectorKafka) Stop() error {
	log.Println("Stop request detected")
	c.stopConsumer()
	c.wg.Wait()
	log.Println("Consuming process stopped")
	if err := c.producer.Close(); err != nil {
		return err
	}
	log.Println("Producer stopped")
	if err := c.consumer.Close(); err != nil {
		return err
	}
	log.Println("Consumer stopped")
	return nil
}
