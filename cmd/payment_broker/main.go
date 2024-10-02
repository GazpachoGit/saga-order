package main

import (
	"bufio"
	"log"
	"os"
	broker "saga-order/internal/broker/payment"
	repo "saga-order/internal/repo/payment"
	service "saga-order/internal/service/payment"
)

func main() {
	repo := repo.NewPaymentRepoV1()
	service := service.NewSagaPaymentServiceV1(repo)
	messageHandler := broker.NewSagaPaymentHandler(service)
	localBroker, err := broker.NewPaymentBroker(messageHandler)
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
