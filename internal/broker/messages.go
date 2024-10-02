package broker

type OrchestratorMessage struct {
	Key       string
	Value     string
	NextTopic string
}

type OrchestratorMessagePayload struct {
	Successes   bool
	Error       string
	AggregateID uint64
	Payload     []byte
}

type CreateOrderMessage struct {
	OrderID    uint64
	ProductID  uint64
	Amount     uint8
	CustomerID uint64
	Cost       uint8
}

type RollbackOrderMessage struct {
	OrderID uint64
}

type CreatePaymentMessage struct {
	OrderID    uint64
	CustomerID uint64
	Amount     uint8
}
