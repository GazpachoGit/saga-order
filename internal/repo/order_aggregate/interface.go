package repo

import "saga-order/internal/model"

type OrderAggregateRepo interface {
	UpdateAggregateState(orderID uint64, state string) error
	GetAggregate(orderID uint64) (model.OrderAggregate, error)
}
