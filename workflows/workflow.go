package workflows

import (
	"eth-temporal/app"
	"eth-temporal/app/activities"
	"time"

	"go.temporal.io/sdk/workflow"
)

func GetLatestBlockNumWorkflow(ctx workflow.Context) (uint64, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx1 := workflow.WithActivityOptions(ctx, options)

	var blockNumber uint64
	err := workflow.ExecuteActivity(ctx1, activities.GetLatestBlockNum).Get(ctx1, &blockNumber)
	if err != nil {
		panic(err)
	}

	return blockNumber, err
}

func GetBlockWorkflow(ctx workflow.Context, blockNumber uint64, backfill bool) (app.Block, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
	}
	ctx1 := workflow.WithActivityOptions(ctx, options)

	var block app.Block
	var traces []app.Trace
	var result string
	var err error
	// Looping to catch updates. Arbitrarily choosing 3 fetches.
	for i := 0; i < 3; i++ {
		// Fetch block
		err = workflow.ExecuteActivity(ctx1, activities.GetBlockByNumber, blockNumber).Get(ctx1, &block)
		if err != nil {
			panic(err)
		}

		// Persist to Postgres
		err = workflow.ExecuteActivity(ctx1, activities.UpsertBlockToPostgres, block).Get(ctx1, &result)
		if err != nil {
			panic(err)
		}

		// Fetch traces
		err = workflow.ExecuteActivity(ctx1, activities.GetTracesByBlock, blockNumber).Get(ctx1, &traces)
		if err != nil {
			panic(err)
		}

		// Persist traces
		err = workflow.ExecuteActivity(ctx1, activities.UpsertTracesToPostgres, blockNumber, traces).Get(ctx1, &result)
		if err != nil {
			panic(err)
		}

		// If it is a backfill, we just break on first iteration
		if backfill {
			break
		}
		workflow.Sleep(ctx1, time.Second*15)
	}
	return block, err
}
