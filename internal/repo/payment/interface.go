package repo

import (
	"saga-order/internal/model"
)

type PaymentRepo interface {
	CreatePayment(createInput model.Payment) error
	DeletePayment(paymentID uint64) error
}
