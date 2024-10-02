package service

import "saga-order/internal/model"

type SagaPaymentService interface {
	CreatePayment(createInput model.Payment) error
	RollbackPayment(paymentID uint64) error
}
