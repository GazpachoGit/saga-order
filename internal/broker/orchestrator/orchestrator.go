package broker

import "saga-order/internal/broker"

const (
	topicWithRequests = "saga_orchestrator"
	consumerGroupName = "saga_order_orchestrator"
)

func NewOrchestratorBroker(h broker.Handler) (broker.Connector, error) {
	br, err := broker.NewBrokerKafka(topicWithRequests, consumerGroupName)
	if err != nil {
		return nil, err
	}
	br.RegisterHandler(h)
	return br, nil
}
