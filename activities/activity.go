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

func ConvertBlock(ctx context.Context, block *web3.Block) (app.Block, error) {
	var transactions []app.Transaction
	for _, t := range block.Transactions {
		transaction := app.Transaction{
			Hash:        t.Hash.String(),
			From:        t.From.String(),
			To:          t.To.String(),
			GasPrice:    t.GasPrice,
			Gas:         t.Gas,
			Value:       t.Value,
			Nonce:       t.Nonce,
			BlockHash:   t.BlockHash.String(),
			BlockNumber: t.BlockNumber,
			TxnIndex:    t.TxnIndex,
		}
		transactions = append(transactions, transaction)
	}

	transactionsJson, err := json.Marshal(transactions)
	if err != nil {
		panic(err)
	}

	newBlock := app.Block{
		Number:           block.Number,
		Hash:             block.Hash.String(),
		ParentHash:       block.ParentHash.String(),
		Sha3Uncles:       block.Sha3Uncles.String(),
		TransactionsRoot: block.TransactionsRoot.String(),
		StateRoot:        block.StateRoot.String(),
		ReceiptsRoot:     block.ReceiptsRoot.String(),
		Miner:            block.Miner.String(),
		Difficulty:       block.Difficulty,
		ExtraData:        string(block.ExtraData),
		GasLimit:         block.GasLimit,
		GasUsed:          block.GasUsed,
		Timestamp:        block.Timestamp,
		Transactions:     string(transactionsJson),
	}

	return newBlock, nil
}

func GetBlockByNumber(ctx context.Context, number uint64) (*web3.Block, error) {
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

	return result, nil
}

func GetLastInsertedBlockNumber(ctx context.Context) (uint64, error) {
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

	rows, err := db.Query("SELECT MAX(number) FROM ethereum.blocks;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var blockNumber uint64
	rows.Next()
	rows.Scan(&blockNumber)
	return blockNumber, err
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
			transactions      JSONB,
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
		UPDATE SET transactions = EXCLUDED.transactions,
		           gas_limit    = EXCLUDED.gas_limit,
				   gas_used     = EXCLUDED.gas_used,
				   timestamp    = EXCLUDED.timestamp,
				   difficulty   = EXCLUDED.difficulty,
				   sha3_uncles  = EXCLUDED.sha3_uncles
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
