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

	var block app.Block
	err = workflow.ExecuteActivity(ctx1, activities.GetBlockByNumber, blockNumber).Get(ctx1, &block)

	if err != nil {
		panic(err)
	}

	var result string
	err = workflow.ExecuteActivity(ctx1, activities.UpsertToPostgres, block).Get(ctx1, &result)

	return block, err
}

func GetBlockWorkflow(ctx workflow.Context, blockNumber uint64) (app.Block, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx1 := workflow.WithActivityOptions(ctx, options)

	var block app.Block
	err := workflow.ExecuteActivity(ctx1, activities.GetBlockByNumber, blockNumber).Get(ctx1, &block)

	if err != nil {
		panic(err)
	}

	// Persist to Postgres
	var result string
	err = workflow.ExecuteActivity(ctx1, activities.UpsertToPostgres, block).Get(ctx1, &result)

	// Poll for updates

	// Upsert to Postgres

	return block, err
}
