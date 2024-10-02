package orchestrator

import (
	"encoding/json"
	broker "saga-order/internal/broker"
	"saga-order/internal/model"
	"saga-order/internal/util"

	"github.com/IBM/sarama"
)

const (
	brokerList        = "localhost:9092"
	maxRetry          = 5
	CREATE_ORDER_CMD  = "CREATE_ORDER_CMD"
	orderServiceTopic = "saga_order_requests"
)

//execute incoming commands

type OrderOrchestratorV1 struct {
	producer sarama.SyncProducer
}

func NewOrderOrchestratorV1() (*OrderOrchestratorV1, error) {
	resp := &OrderOrchestratorV1{}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = maxRetry
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{brokerList}, config)
	if err != nil {
		return nil, err
	}
	resp.producer = producer

	return resp, nil
}

func (orch *OrderOrchestratorV1) Stop() {
	orch.producer.Close()
}

func (orch *OrderOrchestratorV1) StartCreateOrderTransaction(createInput model.Order) error {
	payload := broker.CreateOrderMessage{
		OrderID:    util.Uint64(),
		ProductID:  createInput.ProductID,
		Amount:     createInput.Amount,
		CustomerID: createInput.CustomerID,
		Cost:       createInput.Cost,
	}
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	message := &sarama.ProducerMessage{
		Topic: orderServiceTopic,
		Value: sarama.StringEncoder(payloadStr),
		Key:   sarama.StringEncoder(CREATE_ORDER_CMD),
	}
	_, _, err = orch.producer.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}
