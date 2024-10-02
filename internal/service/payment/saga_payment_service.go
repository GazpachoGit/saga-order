package service

import (
	"saga-order/internal/model"
	repo "saga-order/internal/repo/payment"
)

type SagaPaymentServiceV1 struct {
	repo repo.PaymentRepo
}

func NewSagaPaymentServiceV1(repo repo.PaymentRepo) SagaPaymentService {
	return &SagaPaymentServiceV1{
		repo: repo,
	}
}

func (s *SagaPaymentServiceV1) CreatePayment(createInput model.Payment) error {
	if err := s.repo.CreatePayment(createInput); err != nil {
		return err
	}
	return nil
}

func (s *SagaPaymentServiceV1) RollbackPayment(orderID uint64) error {
	if err := s.repo.DeletePayment(orderID); err != nil {
		return err
	}
	return nil
}
