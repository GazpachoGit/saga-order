package service

import repo "saga-order/internal/repo/order_aggregate"

const (
	ORDER_AGG_ORD_PENDING   = "OrderCreating"
	ORDER_AGG_PAY_PENDING   = "PaymentCreating"
	ORDER_AGG_ORD_CREATED   = "OrderCreated"
	ORDER_AGG_PAY_REJECTING = "PaymentRejecting"
	ORDER_AGG_ORD_REJECTING = "OrderRejecting"
	ORDER_AGG_ORD_REJECTED  = "OrderRejected"
)

type OrderAggregateServiceV1 struct {
	repo repo.OrderAggregateRepo
}

func NewOrderAggregateServiceV1(repo repo.OrderAggregateRepo) OrderAggregateService {
	return &OrderAggregateServiceV1{
		repo: repo,
	}
}

func (s *OrderAggregateServiceV1) CreateAggregate(orderID uint64) error {
	return s.repo.UpdateAggregateState(orderID, ORDER_AGG_ORD_PENDING)
}
func (s *OrderAggregateServiceV1) SetPaymentPending(orderID uint64) error {
	return s.repo.UpdateAggregateState(orderID, ORDER_AGG_PAY_PENDING)
}
func (s *OrderAggregateServiceV1) SetPaymentConfirmed(orderID uint64) error {
	return s.repo.UpdateAggregateState(orderID, ORDER_AGG_ORD_CREATED)
}
func (s *OrderAggregateServiceV1) SetOrderRejecting(orderID uint64) error {
	return s.repo.UpdateAggregateState(orderID, ORDER_AGG_ORD_REJECTING)
}
func (s *OrderAggregateServiceV1) SetPaymentRejecting(orderID uint64) error {
	return s.repo.UpdateAggregateState(orderID, ORDER_AGG_PAY_REJECTING)
}
func (s *OrderAggregateServiceV1) SetOrderRejected(orderID uint64) error {
	return s.repo.UpdateAggregateState(orderID, ORDER_AGG_ORD_REJECTED)
}
