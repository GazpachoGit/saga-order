package repo

import (
	"saga-order/internal/model"
)

type PaymentRepoV1 struct {
	db map[uint64]model.Payment
}

func NewPaymentRepoV1() PaymentRepo {
	return &PaymentRepoV1{
		db: make(map[uint64]model.Payment, 5),
	}
}

func (repo *PaymentRepoV1) CreatePayment(createInput model.Payment) error {
	if repo.db == nil {
		repo.db = make(map[uint64]model.Payment)
	}
	//TODO generate ID
	repo.db[createInput.ID] = createInput
	return nil
}

func (repo *PaymentRepoV1) DeletePayment(paymentID uint64) error {
	if repo.db == nil {
		repo.db = make(map[uint64]model.Payment)
		return nil
	}
	delete(repo.db, paymentID)
	return nil
}
