package workflows

import (
	"eth-temporal/app"
	"eth-temporal/app/activities"
	"time"

	"github.com/umbracle/go-web3"
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

func GetLatestBlockWorkflow(ctx workflow.Context) (app.Block, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx1 := workflow.WithActivityOptions(ctx, options)

	var blockNumber uint64
	err := workflow.ExecuteActivity(ctx1, activities.GetLatestBlockNum).Get(ctx1, &blockNumber)
	if err != nil {
		panic(err)
	}

	var rawBlock web3.Block
	err = workflow.ExecuteActivity(ctx1, activities.GetBlockByNumber, blockNumber).Get(ctx1, &rawBlock)
	if err != nil {
		panic(err)
	}

	var block app.Block
	err = workflow.ExecuteActivity(ctx1, activities.ConvertBlock, rawBlock).Get(ctx1, &block)
	if err != nil {
		panic(err)
	}

	var result string
	err = workflow.ExecuteActivity(ctx1, activities.UpsertToPostgres, block).Get(ctx1, &result)
	if err != nil {
		panic(err)
	}

	return block, err
}

func GetBlockWorkflow(ctx workflow.Context, blockNumber uint64) (app.Block, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
	}
	ctx1 := workflow.WithActivityOptions(ctx, options)

	var block app.Block
	var result string
	var err error
	// Looping to catch updates. Arbitrarily choosing 5 loops.
	for i := 0; i < 5; i++ {
		// Fetch block
		err = workflow.ExecuteActivity(ctx1, activities.GetBlockByNumber, blockNumber).Get(ctx1, &block)
		if err != nil {
			panic(err)
		}
		// Persist to Postgres
		err = workflow.ExecuteActivity(ctx1, activities.UpsertToPostgres, block).Get(ctx1, &result)
		if err != nil {
			panic(err)
		}
		workflow.Sleep(ctx1, time.Second*15)
	}
	return block, err
}
