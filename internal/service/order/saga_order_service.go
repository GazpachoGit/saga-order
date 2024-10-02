package service

import (
	"saga-order/internal/model"
	repo "saga-order/internal/repo/order"
)

type SagaOrderServiceV1 struct {
	repo repo.OrderRepo
}

func NewSagaOrderServiceV1(repo repo.OrderRepo) SagaOrderService {
	return &SagaOrderServiceV1{
		repo: repo,
	}
}

func (s *SagaOrderServiceV1) CreateOrder(createInput model.Order) (model.Order, error) {
	o, err := s.repo.CreateOrder(createInput)
	if err != nil {
		return model.Order{}, err
	}
	return o, nil
}

func (s *SagaOrderServiceV1) RollbackOrder(orderID uint64) error {
	if err := s.repo.DeleteOrder(orderID); err != nil {
		return err
	}
	return nil
}
