package broker

import "saga-order/internal/broker"

const (
	topicWithRequests = "saga_order_requests"
	consumerGroupName = "saga_order"
)

func NewOrderBroker(h broker.Handler) (broker.Connector, error) {
	br, err := broker.NewBrokerKafka(topicWithRequests, consumerGroupName)
	if err != nil {
		return nil, err
	}
	br.RegisterHandler(h)
	return br, nil
}
