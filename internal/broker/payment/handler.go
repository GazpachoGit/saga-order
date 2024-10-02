package broker

import (
	"encoding/json"
	"fmt"
	broker "saga-order/internal/broker"
	"saga-order/internal/model"
	service "saga-order/internal/service/payment"
)

type SagaPaymentHandler struct {
	svc service.SagaPaymentService
}

const (
	CREATE_PAYMENT_CMD        = "CREATE_PAYMENT_CMD"
	CREATE_PAYMENT_RESP_KEY   = "CREATE_PAYMENT_RESP"
	ROLLBACK_PAYMENT          = "ROLLBACK_PAYMENT"
	ROLLBACK_PAYMENT_RESP_KEY = "ROLLBACK_PAYMENT_RESP"
	RESPONSE_TOPIC            = "saga_orchestrator"
)

func NewSagaPaymentHandler(aggregateSrv service.SagaPaymentService) broker.Handler {
	return &SagaPaymentHandler{
		svc: aggregateSrv,
	}
}

func (h *SagaPaymentHandler) Handle(requestKey []byte, requestMsg []byte) (respMsg broker.OrchestratorMessage, err error) {
	strKey := string(requestKey)
	switch strKey {
	case CREATE_PAYMENT_CMD:
		return h.CreatePayment(requestMsg)
	}
	return broker.OrchestratorMessage{}, fmt.Errorf("Unknown message key", strKey)
}

func (h *SagaPaymentHandler) CreatePayment(requestMsg []byte) (respMsg broker.OrchestratorMessage, err error) {
	msg := broker.OrchestratorMessagePayload{}
	createPaymentPayload := broker.CreatePaymentMessage{}
	err = json.Unmarshal(requestMsg, &createPaymentPayload)
	if err != nil {
		return broker.OrchestratorMessage{}, nil
	} else {
		err = h.svc.CreatePayment(model.Payment{
			CustomerID: createPaymentPayload.CustomerID,
			Amount:     createPaymentPayload.Amount,
		})
		if err != nil {
			msg.Successes = false
			msg.Error = err.Error()
			msg.AggregateID = createPaymentPayload.OrderID

		} else {
			msg.Successes = true
			msg.Error = ""
			msg.AggregateID = createPaymentPayload.OrderID
		}

	}
	m, err := json.Marshal(msg)
	if err != nil {
		return broker.OrchestratorMessage{}, err
	}
	return broker.OrchestratorMessage{
		Key:       CREATE_PAYMENT_RESP_KEY,
		Value:     string(m),
		NextTopic: RESPONSE_TOPIC,
	}, nil
}
