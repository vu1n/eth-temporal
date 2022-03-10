package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"eth-temporal/app"
	"eth-temporal/app/workflows"
)

type fetchState struct {
	blockNumber uint64
}

func (f *fetchState) Pop() uint64 {
	result := f.blockNumber
	f.blockNumber++
	return result
}

func startWorkflow(c client.Client, state *fetchState) {
	for {
		blockNum := state.Pop()
		fmt.Printf("Queuing block number %v\n", blockNum)
		wfOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("get-block-%v", blockNum),
			TaskQueue: app.NewBlockTaskQueue,
		}

		wf, err := c.ExecuteWorkflow(context.Background(), wfOptions, workflows.GetBlockWorkflow, blockNum)
		if err != nil {
			fmt.Println(err)
			log.Fatalln("unable to get block 1")
		}
		var block app.Block
		err = wf.Get(context.Background(), &block)
		if err != nil {
			fmt.Println(err)
			log.Fatalln("unable to get block 2")
		}
	}
}

func main() {
	size := flag.Int("size", 1, "The queue size")
	flag.Parse()

	// Create the client object just once per process
	c, err := app.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	options := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("new-block-workflow-%v", time.Now().Unix()),
		TaskQueue: app.NewBlockTaskQueue,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.GetLatestBlockNumWorkflow)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}
	var blockNum uint64
	err = we.Get(context.Background(), &blockNum)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	fmt.Println("Fetching latest blocks")

	currentState := &fetchState{blockNum}
	for i := 0; i < *size; i++ {
		go startWorkflow(c, currentState)
	}
	for {
		time.Sleep(time.Second * 15)
		fmt.Println("Working . . . ")
	}

}
