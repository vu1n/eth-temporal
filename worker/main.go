package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"eth-temporal/app"
	"eth-temporal/app/activities"
	"eth-temporal/app/workflows"
)

func main() {
	// Create the client object just once per process

	c, err := app.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal clientttttt", os.Getenv("TEMPORAL_GRPC_ENDPOINT"), err)
	}
	defer c.Close()
	//This worker hosts both Workflow and Activity functions

	w := worker.New(c, app.NewBlockTaskQueue, worker.Options{})

	w.RegisterActivity(activities.ConvertBlock)
	w.RegisterActivity(activities.GetBlockByNumber)
	w.RegisterActivity(activities.GetLatestBlockNum)
	w.RegisterActivity(activities.UpsertToPostgres)

	w.RegisterWorkflow(workflows.GetLatestBlockNumWorkflow)
	w.RegisterWorkflow(workflows.GetBlockWorkflow)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
