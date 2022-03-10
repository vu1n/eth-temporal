package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"eth-temporal/app"
	"eth-temporal/app/workflows"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.temporal.io/sdk/client"
)

type handlers struct {
	temporalClient client.Client
}

func (h *handlers) handleGetBlockByNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockNumber := vars["blockNumber"]

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", app.DbHost, app.DbPort, app.DbUser, app.DbPassword, app.DbName)

	// Connect to pg db
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// clean up db connection
	defer db.Close()
	fmt.Printf("Fetching %s\n", blockNumber)

	selectSql := fmt.Sprintf(
		`SELECT json_build_object(
			'number', number,
			'hash', hash,
			'parent_hash', parent_hash,
			'sha3_uncles', sha3_uncles,
			'transactions_root', transactions_root,
			'state_root', state_root,
			'receipts_root', receipts_root,
			'miner', miner,
			'difficulty', difficulty,
			'extra_data', extra_data,
			'gas_limit', gas_limit,
			'gas_used', gas_limit,
			'timestamp', timestamp,
			'transactions', transactions
			) block
		FROM ethereum.blocks WHERE number = %v;`, blockNumber)
	fmt.Println(selectSql)
	rows, err := db.Query(selectSql)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var block string
	rows.Next()
	err = rows.Scan(&block)
	if err != nil {
		// On error we will queue a task to fetch the block and return the reuslt
		var block app.Block
		var blockNum uint64
		blockNum, _ = strconv.ParseUint(blockNumber, 10, 64)
		we, err := h.temporalClient.ExecuteWorkflow(
			r.Context(),
			client.StartWorkflowOptions{
				TaskQueue: app.NewBlockTaskQueue,
				ID:        fmt.Sprintf("get-block-from-api-call-%v", blockNumber),
			},
			workflows.GetBlockWorkflow,
			blockNum,
		)
		if err != nil {
			fmt.Printf("failed to start workflow: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		we.Get(context.Background(), &block)
		if err != nil {
			fmt.Printf("failed to retrieve block: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		blockJson, err := json.Marshal(block)
		if err != nil {
			fmt.Printf("failed to marshal json: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Println(blockJson)
		w.Write(blockJson)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println(block)
	w.Write([]byte(block))
}

func Router(c client.Client) *mux.Router {
	r := mux.NewRouter()

	h := handlers{temporalClient: c}

	r.HandleFunc("/blockNumber/{blockNumber:[0-9]+}", h.handleGetBlockByNumber).Methods("GET")

	return r
}
