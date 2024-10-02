package main

import (
	"bufio"
	"log"
	"os"
	broker "saga-order/internal/broker/orchestrator"
	repo "saga-order/internal/repo/order_aggregate"
	service "saga-order/internal/service/order_aggregate"
)

func main() {
	repo := repo.NewOrderAggregateRepoV1()
	service := service.NewOrderAggregateServiceV1(repo)
	messageHandler := broker.NewSagaOrchestratorHandler(service)
	localBroker, err := broker.NewOrchestratorBroker(messageHandler)
	if err != nil {
		panic(err)
	}
	log.Println("Starting broker...")
	go localBroker.Run()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		localBroker.Stop()
		break
	}
}
