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

func startWorkflow(c client.Client, state *fetchState, end uint64) {
	for {
		if state.blockNumber > end {
			break
		}
		blockNum := state.Pop()
		fmt.Printf("Queuing block number %v\n", blockNum)
		wfOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("get-block-%v", blockNum),
			TaskQueue: app.BackfillTaskQueue,
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
	size := flag.Int("size", 2, "The queue size")
	start := flag.Uint64("start", 1, "The starting block number")
	end := flag.Uint64("end", 1, "The ending block number")
	flag.Parse()

	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	fmt.Printf("Fetching blocks from %v to %v\n", start, end)

	currentState := &fetchState{*start}
	for i := 0; i < *size; i++ {
		go startWorkflow(c, currentState, *end)
	}
	for {
		if currentState.blockNumber > *end {
			break
		}
		time.Sleep(time.Second * 15)
		fmt.Println("Working . . . ")
	}

}
