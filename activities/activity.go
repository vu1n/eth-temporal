package activities

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"eth-temporal/app"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lib/pq"
	web3 "github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"go.temporal.io/sdk/activity"
)

func GetLatestBlockNum(ctx context.Context) (uint64, error) {
	logger := activity.GetLogger(ctx)

	client, err := jsonrpc.NewClient(app.RpcHost)
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

	client, err := jsonrpc.NewClient(app.RpcHost)
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("Fetching block: %v\n", number))

	result, err := client.Eth().GetBlockByNumber(web3.BlockNumber(number), true)
	if err != nil {
		panic(err)
	}

	if result == nil {
		panic("No results")
	}

	var transactions []app.Transaction
	for _, t := range result.Transactions {
		transaction := app.Transaction{
			Hash:        t.Hash.String(),
			From:        t.From.String(),
			To:          "",
			GasPrice:    t.GasPrice,
			Gas:         t.Gas,
			Value:       t.Value,
			Nonce:       t.Nonce,
			BlockHash:   t.BlockHash.String(),
			BlockNumber: t.BlockNumber,
			TxnIndex:    t.TxnIndex,
		}
		if t.To != nil {
			transaction.To = t.To.String()
		}
		transactions = append(transactions, transaction)
	}

	transactionsJson, err := json.Marshal(transactions)
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
		ExtraData:        string(bytes.Split(result.ExtraData[:], []byte{0})[0]),
		GasLimit:         result.GasLimit,
		GasUsed:          result.GasUsed,
		Timestamp:        result.Timestamp,
		Transactions:     string(transactionsJson),
	}

	return block, nil
}

func GetTracesByBlock(ctx context.Context, number uint64) ([]app.Trace, error) {
	logger := activity.GetLogger(ctx)

	logger.Info(fmt.Sprintf("Fetching traces for block: %v\n", number))

	postBody, _ := json.Marshal(app.TraceBlockPayload{
		Method:  "trace_block",
		Params:  []string{fmt.Sprintf("0x%x", number)},
		Id:      1,
		Jsonrpc: "2.0",
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("https://eth-rpc.gateway.pokt.network", "application/json", responseBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var res app.TraceBlockResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	if res.Error != nil {
		panic(res.Error)
	}

	if len(res.Result) == 0 {
		panic("No results")
	}
	return res.Result, nil
}

func GetLastInsertedBlockNumber(ctx context.Context) (uint64, error) {
	logger := activity.GetLogger(ctx)
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", app.DbHost, app.DbPort, app.DbUser, app.DbPassword, app.DbName)

	// Connect to pg db
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	logger.Info(fmt.Sprintf("Connected to %s", app.DbHost))
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

func UpsertBlockToPostgres(ctx context.Context, block app.Block) error {
	logger := activity.GetLogger(ctx)
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", app.DbHost, app.DbPort, app.DbUser, app.DbPassword, app.DbName)

	// Connect to pg db
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	logger.Info(fmt.Sprintf("Connected to %s", app.DbHost))
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
			gas_limit         BIGINT         DEFAULT NULL,
			gas_used          BIGINT         DEFAULT NULL,
			timestamp         BIGINT         NOT NULL,
			transactions      JSON,
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
		block.Miner, block.Difficulty, block.GasLimit, block.GasUsed, block.Timestamp, block.Transactions)
	// logger.Info("Executing:")
	// fmt.Println(upsertSql)
	_, err = db.Exec(upsertSql)
	if err != nil {
		fmt.Println(err)
		fmt.Println(upsertSql)
		panic(err)
	}

	return nil
}

func UpsertTracesToPostgres(ctx context.Context, blockNumber uint64, traces []app.Trace) error {
	logger := activity.GetLogger(ctx)
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", app.DbHost, app.DbPort, app.DbUser, app.DbPassword, app.DbName)

	// Connect to pg db
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	logger.Info(fmt.Sprintf("Connected to %s", app.DbHost))
	// clean up db connection
	defer db.Close()

	txn, err := db.Begin()
	if err != nil {
		panic(err)
	}

	// create table
	createSql :=
		`CREATE SCHEMA IF NOT EXISTS ethereum;
		CREATE TABLE IF NOT EXISTS ethereum.traces (
			block_number      BIGINT         NOT NULL,
			block_hash        TEXT           NOT NULL,
			transaction_hash  TEXT           NOT NULL,
			from_address      TEXT           NOT NULL,
			to_address        TEXT           NOT NULL,
			value             NUMERIC(38,0)  NOT NULL,
			input             TEXT           NOT NULL,
			output            TEXT           NOT NULL,
			trace_type        TEXT           NOT NULL,
			call_type         TEXT           NOT NULL,
			reward_type       TEXT           NOT NULL,
			gas               BIGINT         NOT NULL,
			gas_used          BIGINT         NOT NULL,
			subtraces         BIGINT         NOT NULL,
			trace_address     TEXT           NOT NULL,
			transaction_pos   BIGINT         NOT NULL,
			error             TEXT           NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_traces_block_number ON ethereum.traces (block_number);
		CREATE TEMP TABLE traces_stage AS TABLE ethereum.traces WITH NO DATA;
		`

	_, err = txn.Exec(createSql)
	if err != nil {
		panic(err)
	}

	// Copy data into temp table for upsert
	stmt, err := txn.Prepare(pq.CopyIn(
		"traces_stage",
		"block_number", "block_hash", "transaction_hash", "from_address", "to_address", "value",
		"input", "output", "trace_type", "call_type", "reward_type", "gas",
		"gas_used", "subtraces", "trace_address", "transaction_pos", "error"))
	if err != nil {
		panic(err)
	}

	for _, trace := range traces {
		traceAddress, _ := json.Marshal(trace.TraceAddress)
		_, err = stmt.Exec(trace.BlockNumber, trace.BlockHash, trace.TransactionHash,
			trace.Action.From, trace.Action.To, app.HexToFloat(trace.Action.Value), trace.Action.Input,
			trace.Result.Output, trace.Type, trace.Action.CallType, trace.Action.RewardType, app.HexToUInt(trace.Action.Gas),
			app.HexToUInt(trace.Result.GasUsed), trace.Subtraces, traceAddress, trace.TransactionPosition, trace.Error)
		if err != nil {
			panic(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	err = stmt.Close()
	if err != nil {
		panic(err)
	}

	// upsert traces
	// I'm not sure if we can UPSERT, so will just DELETE rows and insert
	upsertSql := fmt.Sprintf(`
		DELETE FROM ethereum.traces WHERE block_number = '%v';
		INSERT INTO ethereum.traces
		SELECT * FROM traces_stage
		`, blockNumber)
	// ON CONFLICT (block_number, transaction_hash, transaction_pos, from_address, to_address, trace_address)
	// DO
	// UPDATE SET
	// 		value         = EXCLUDED.value,
	// 		input         = EXCLUDED.input,
	// 		output        = EXCLUDED.output,
	// 		trace_type    = EXCLUDED.trace_type,
	// 		call_type     = EXCLUDED.call_type,
	// 		reward_type   = EXCLUDED.reward_type,
	// 		gas           = EXCLUDED.gas,
	// 		gas_used      = EXCLUDED.gas_used,
	// 		subtraces     = EXCLUDED.subtraces,
	// 		error         = EXCLUDED.error
	_, err = txn.Exec(upsertSql)
	if err != nil {
		fmt.Println(err)
		fmt.Println(upsertSql)
		panic(err)
	}

	err = txn.Commit()
	if err != nil {
		panic(err)
	}

	return nil
}
