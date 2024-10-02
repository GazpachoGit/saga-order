package service

type OrderAggregateService interface {
	CreateAggregate(orderID uint64) error
	SetPaymentPending(orderID uint64) error
	SetPaymentConfirmed(orderID uint64) error
	SetOrderRejecting(orderID uint64) error
	SetOrderRejected(orderID uint64) error
}
