package model

type Order struct {
	ID         uint64
	ProductID  uint64
	Amount     uint8
	CustomerID uint64
	Cost       uint8
	//UpdatedAt  int64
	//CreatedAt  int64
}
