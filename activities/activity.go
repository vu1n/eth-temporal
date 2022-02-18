package activities

import (
	"context"
	"database/sql"
	"encoding/json"
	"eth-temporal/app"
	"fmt"

	_ "github.com/lib/pq"
	web3 "github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"go.temporal.io/sdk/activity"
)

const (
	rpcHost  = "https://eth-rpc.gateway.pokt.network"
	host     = "localhost"
	port     = 5433
	user     = "temporal"
	password = "temporal"
	dbname   = "postgres"
)

func GetLatestBlockNum(ctx context.Context) (uint64, error) {
	logger := activity.GetLogger(ctx)

	client, err := jsonrpc.NewClient(rpcHost)
	if err != nil {
		return 0, err
	}

	logger.Info("\nFetching latest block number\n")

	number, err := client.Eth().BlockNumber()
	if err != nil {
		return 0, err
	}

	logger.Info(fmt.Sprintf("\nLatest block number: %v\n", number))

	return number, err
}

func GetBlockByNumber(ctx context.Context, number uint64) (app.Block, error) {
	logger := activity.GetLogger(ctx)

	client, err := jsonrpc.NewClient(rpcHost)
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("Fetching block: %v\n", number))

	result, err := client.Eth().GetBlockByNumber(web3.BlockNumber(number), true)
	if err != nil {
		panic(err)
	}

	transactions, err := json.Marshal(result.Transactions)
	if err != nil {
		panic(err)
	}

	block := app.Block{
		Number:           result.Number,
		Hash:             result.Hash.String(),
		ParentHash:       result.ParentHash.String(),
		Sha3Uncles:       result.Sha3Uncles.String(),
		TransactionsRoot: result.TransactionsRoot.String(),
		StateRoot:        result.StateRoot.String(),
		ReceiptsRoot:     result.ReceiptsRoot.String(),
		Miner:            result.Miner.String(),
		Difficulty:       result.Difficulty,
		ExtraData:        string(result.ExtraData),
		GasLimit:         result.GasLimit,
		GasUsed:          result.GasUsed,
		Timestamp:        result.Timestamp,
		Transactions:     string(transactions),
	}

	return block, nil
}

func UpsertToPostgres(ctx context.Context, block app.Block) error {
	logger := activity.GetLogger(ctx)
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Connect to pg db
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	logger.Info(fmt.Sprintf("Connected to %s", host))
	// clean up db connection
	defer db.Close()

	// create table
	createSql :=
		`CREATE SCHEMA IF NOT EXISTS ethereum;
		CREATE TABLE IF NOT EXISTS ethereum.blocks (
			number            BIGINT         NOT NULL,
			hash              TEXT           NOT NULL,
			parent_hash       TEXT           NOT NULL,
			sha3_uncles       TEXT           NOT NULL,
			transactions_root TEXT           NOT NULL,
			state_root        TEXT           NOT NULL,
			receipts_root     TEXT           NOT NULL,
			miner             TEXT           NOT NULL,
			difficulty        NUMERIC(38, 0) NOT NULL,
			extra_data        TEXT           DEFAULT NULL,
			gas_limit         BIGINT         DEFAULT NULL,
			gas_used          BIGINT         DEFAULT NULL,
			timestamp         BIGINT         NOT NULL,
			transactions      TEXT,
			PRIMARY KEY (number)
		)`

	_, err = db.Exec(createSql)
	if err != nil {
		panic(err)
	}
	// upsert block
	upsertSql := fmt.Sprintf(
		`INSERT INTO ethereum.blocks (
			number,
			hash,
			parent_hash,
			sha3_uncles,
			transactions_root,
			state_root,
			receipts_root,
			miner,
			difficulty,
			extra_data,
			gas_limit,
			gas_used,
			timestamp,
			transactions
		) VALUES (
			'%v', 
			'%v', 
			'%v', 
			'%v', 
			'%v', 
			'%v', 
			'%v', 
			'%v', 
			'%v', 
			'%v', 
			'%v', 
			'%v', 
			'%v', 
			'%v'
		)
		ON CONFLICT(number)
		DO
		UPDATE SET transactions = EXCLUDED.transactions
		`, block.Number, block.Hash, block.ParentHash, block.Sha3Uncles, block.TransactionsRoot, block.StateRoot, block.ReceiptsRoot,
		block.Miner, block.Difficulty, block.ExtraData, block.GasLimit, block.GasUsed, block.Timestamp, block.Transactions)
	logger.Info("Executing:")
	fmt.Println(upsertSql)
	_, err = db.Exec(upsertSql)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return err
}
