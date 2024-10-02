package repo

import (
	"fmt"
	"log"
	"saga-order/internal/model"
)

type OrderAggregateRepoV1 struct {
	db map[uint64]model.OrderAggregate
}

func NewOrderAggregateRepoV1() OrderAggregateRepo {
	return &OrderAggregateRepoV1{
		db: make(map[uint64]model.OrderAggregate, 5),
	}
}

func (repo *OrderAggregateRepoV1) UpdateAggregateState(orderID uint64, state string) error {
	if repo.db == nil {
		repo.db = make(map[uint64]model.OrderAggregate)
	}
	repo.db[orderID] = model.OrderAggregate{
		ID:    orderID,
		State: state,
	}
	log.Printf("Aggregate repo. Update state. ID: %v, State: %s", orderID, state)
	return nil
}

func (repo *OrderAggregateRepoV1) GetAggregate(orderID uint64) (model.OrderAggregate, error) {
	if repo.db == nil {
		repo.db = make(map[uint64]model.OrderAggregate)
	}
	value, ok := repo.db[orderID]
	if !ok {
		return model.OrderAggregate{}, fmt.Errorf("object not found in the repo")
	}
	return value, nil
}
