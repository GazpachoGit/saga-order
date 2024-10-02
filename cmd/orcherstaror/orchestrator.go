package main

import (
	"bufio"
	"log"
	"os"
	"saga-order/internal/model"
	"saga-order/internal/orchestrator"
)

func main() {
	orch, err := orchestrator.NewOrderOrchestratorV1()
	if err != nil {
		panic(err)
	}
	defer orch.Stop()

	scanner := bufio.NewScanner(os.Stdin)
	var i uint64 = 1
	for scanner.Scan() {
		log.Printf("Start transaction with i = %v", i)
		createInput := model.Order{
			ProductID:  i,
			Amount:     uint8(i),
			CustomerID: i,
			Cost:       uint8(i),
		}
		err := orch.StartCreateOrderTransaction(createInput)
		if err != nil {
			log.Println(err)
		}
		i++
	}
}
