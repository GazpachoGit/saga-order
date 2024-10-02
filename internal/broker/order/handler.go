package broker

import (
	"encoding/json"
	"saga-order/internal/broker"
	"saga-order/internal/model"
	service "saga-order/internal/service/order"
)

type SagaOrderHandler struct {
	svc service.SagaOrderService
}

const (
	CREATE_ORDER_CMD        = "CREATE_ORDER_CMD"
	CREATE_ORDER_RESP_KEY   = "CREATE_ORDER_RESP"
	ROLLBACK_ORDER          = "ROLLBACK_ORDER"
	ROLLBACK_ORDER_RESP_KEY = "ROLLBACK_ORDER_RESP"
	RESPONSE_TOPIC          = "saga_orchestrator"
)

func NewSagaOrderHandler(aggregateSrv service.SagaOrderService) broker.Handler {
	return &SagaOrderHandler{
		svc: aggregateSrv,
	}
}

func (h *SagaOrderHandler) Handle(requestKey []byte, requestMsg []byte) (respMsg broker.OrchestratorMessage, err error) {
	strKey := string(requestKey)
	switch strKey {
	case CREATE_ORDER_CMD:
		return h.CreateOrder(requestMsg)
	case ROLLBACK_ORDER:
		return h.RollbackOrder(requestMsg)
	}

	return broker.OrchestratorMessage{}, nil
}

func (h *SagaOrderHandler) CreateOrder(requestMsg []byte) (respMsg broker.OrchestratorMessage, err error) {
	msg := broker.OrchestratorMessagePayload{}
	createOrderMsg := broker.CreateOrderMessage{}
	err = json.Unmarshal(requestMsg, &createOrderMsg)
	if err != nil {
		return broker.OrchestratorMessage{}, err
	} else {
		msg.AggregateID = createOrderMsg.OrderID
		order, err := h.svc.CreateOrder(model.Order{
			ID:         createOrderMsg.OrderID,
			ProductID:  createOrderMsg.ProductID,
			Amount:     createOrderMsg.Amount,
			CustomerID: createOrderMsg.CustomerID,
			Cost:       createOrderMsg.Cost,
		})
		if err != nil {
			msg.Successes = false
			msg.Error = err.Error()
			msg.AggregateID = createOrderMsg.OrderID

		}
		payload, err := json.Marshal(order)
		if err != nil {
			return broker.OrchestratorMessage{}, err
		}
		msg.Successes = true
		msg.Error = ""
		msg.Payload = payload
		msg.AggregateID = createOrderMsg.OrderID

		m, err := json.Marshal(msg)
		if err != nil {
			return broker.OrchestratorMessage{}, err
		}
		return broker.OrchestratorMessage{
			Key:       CREATE_ORDER_RESP_KEY,
			Value:     string(m),
			NextTopic: RESPONSE_TOPIC,
		}, nil
	}

}

func (h *SagaOrderHandler) RollbackOrder(requestMsg []byte) (respMsg broker.OrchestratorMessage, err error) {
	msg := broker.OrchestratorMessagePayload{}
	rollbackOrderMsg := broker.RollbackOrderMessage{}
	err = json.Unmarshal(requestMsg, &rollbackOrderMsg)
	if err != nil {
		return broker.OrchestratorMessage{}, err
	}
	err = h.svc.RollbackOrder(rollbackOrderMsg.OrderID)
	if err != nil {
		msg.Successes = false
		msg.Error = err.Error()
		msg.AggregateID = rollbackOrderMsg.OrderID

	} else {
		msg.Successes = true
		msg.Error = ""
		msg.AggregateID = rollbackOrderMsg.OrderID
	}

	m, err := json.Marshal(msg)
	if err != nil {
		return broker.OrchestratorMessage{}, err
	}
	return broker.OrchestratorMessage{
		Key:       ROLLBACK_ORDER_RESP_KEY,
		Value:     string(m),
		NextTopic: RESPONSE_TOPIC,
	}, nil
}
