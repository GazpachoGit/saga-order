package repo

import (
	"saga-order/internal/model"
)

type OrderRepoV1 struct {
	db map[uint64]model.Order
}

func NewOrderRepoV1() OrderRepo {
	return &OrderRepoV1{
		db: make(map[uint64]model.Order, 5),
	}
}

func (repo *OrderRepoV1) CreateOrder(createInput model.Order) (model.Order, error) {
	if repo.db == nil {
		repo.db = make(map[uint64]model.Order)
	}
	repo.db[createInput.ID] = createInput
	return createInput, nil
}

func (repo *OrderRepoV1) DeleteOrder(orderID uint64) error {
	if repo.db == nil {
		repo.db = make(map[uint64]model.Order)
		return nil
	}
	delete(repo.db, orderID)
	return nil
}
