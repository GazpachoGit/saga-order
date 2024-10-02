package main

import (
	"bufio"
	"log"
	"os"
	broker "saga-order/internal/broker/order"
	repo "saga-order/internal/repo/order"
	service "saga-order/internal/service/order"
)

func main() {
	repo := repo.NewOrderRepoV1()
	service := service.NewSagaOrderServiceV1(repo)
	messageHandler := broker.NewSagaOrderHandler(service)
	localBroker, err := broker.NewOrderBroker(messageHandler)
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
