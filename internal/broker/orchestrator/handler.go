package broker

import (
	"encoding/json"
	"fmt"
	"saga-order/internal/broker"
	"saga-order/internal/model"
	service "saga-order/internal/service/order_aggregate"
)

const (
	CREATE_ORDER_CMD          = "START_ORDER_CREATION"
	CREATE_PAYMENT_CMD        = "CREATE_PAYMENT_CMD"
	CREATE_ORDER_RESP_KEY     = "CREATE_ORDER_RESP"
	ROLLBACK_ORDER_RESP_KEY   = "ROLLBACK_ORDER_RESP"
	CREATE_PAYMENT_RESP_KEY   = "CREATE_PAYMENT_RESP"
	ROLLBACK_PAYMENT_RESP_KEY = "ROLLBACK_PAYMENT_RESP"

	PAYMENT_TOPIC = "saga_payment_requests"
	ORDER_TOPIC   = "saga_order_requests"
)

type SagaOrchestratorHandler struct {
	aggregateSrv service.OrderAggregateService
}

func NewSagaOrchestratorHandler(aggregateSrv service.OrderAggregateService) broker.Handler {
	return &SagaOrchestratorHandler{
		aggregateSrv: aggregateSrv,
	}
}

func (orch *SagaOrchestratorHandler) Handle(requestKey []byte, requestMsg []byte) (respMsg broker.OrchestratorMessage, err error) {
	strKey := string(requestKey)
	switch strKey {
	case CREATE_ORDER_RESP_KEY:
		return orch.ProcessOrderCreationResult(requestMsg)
	case CREATE_PAYMENT_RESP_KEY:
		//finish or start the order rollback
		return orch.ProcessPaymentCreationResult(requestMsg)
	case ROLLBACK_ORDER_RESP_KEY:
		return orch.ProcessOrderCreationResult(requestMsg)
	}

	return broker.OrchestratorMessage{}, fmt.Errorf("Unknown message key")
}

func (h *SagaOrchestratorHandler) ProcessOrderCreationResult(orderCreateMsg []byte) (respMsg broker.OrchestratorMessage, err error) {
	msg := broker.OrchestratorMessagePayload{}
	err = json.Unmarshal(orderCreateMsg, &msg)
	if err != nil {
		return broker.OrchestratorMessage{}, err
	}
	if msg.Successes == true {
		createdOrder := model.Order{}
		err = json.Unmarshal(msg.Payload, &createdOrder)
		if err != nil {
			return broker.OrchestratorMessage{}, err
		}
		nextMsg := broker.CreatePaymentMessage{
			OrderID:    createdOrder.ID,
			CustomerID: createdOrder.CustomerID,
			Amount:     createdOrder.Cost,
		}
		m, err := json.Marshal(nextMsg)
		if err != nil {
			return broker.OrchestratorMessage{}, err
		}
		h.aggregateSrv.SetPaymentPending(msg.AggregateID)
		return broker.OrchestratorMessage{
			Key:       CREATE_PAYMENT_CMD,
			Value:     string(m),
			NextTopic: PAYMENT_TOPIC,
		}, nil
	} else {
		if err := h.aggregateSrv.SetOrderRejected(msg.AggregateID); err != nil {
			return broker.OrchestratorMessage{}, nil
		}

		return broker.OrchestratorMessage{}, nil
	}
}

func (h *SagaOrchestratorHandler) ProcessPaymentCreationResult(paymentCreateMsg []byte) (respMsg broker.OrchestratorMessage, err error) {
	msg := broker.OrchestratorMessagePayload{}
	err = json.Unmarshal(paymentCreateMsg, &msg)
	if err != nil {
		return broker.OrchestratorMessage{}, err
	}
	if msg.Successes == true {
		h.aggregateSrv.SetPaymentConfirmed(msg.AggregateID)
		return broker.OrchestratorMessage{}, nil
	} else {
		h.aggregateSrv.SetOrderRejecting(msg.AggregateID)
		nextMsg := broker.RollbackOrderMessage{
			OrderID: msg.AggregateID,
		}
		m, err := json.Marshal(nextMsg)
		if err != nil {
			return broker.OrchestratorMessage{}, err
		}
		return broker.OrchestratorMessage{
			Key:       ROLLBACK_ORDER_RESP_KEY,
			Value:     string(m),
			NextTopic: ORDER_TOPIC,
		}, nil
	}
}
