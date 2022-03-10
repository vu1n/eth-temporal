package app

import (
	"math/big"
	"os"

	"go.temporal.io/sdk/client"
)

const NewBlockTaskQueue = "NEW_BLOCK_TASK_QUEUE"
const BackfillTaskQueue = "BACKFILL_TASK_QUEUE"

const RpcHost = "https://eth-rpc.gateway.pokt.network"
const DbHost = "eth-pg"
const DbPort = 5433
const DbUser = "temporal"
const DbPassword = "temporal"
const DbName = "postgres"

type Transaction struct {
	Hash        string
	From        string
	To          string
	GasPrice    uint64
	Gas         uint64
	Value       *big.Int
	Nonce       uint64
	BlockHash   string
	BlockNumber uint64
	TxnIndex    uint64
}

type Block struct {
	Number           uint64   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parent_hash"`
	Sha3Uncles       string   `json:"sha3_uncles"`
	TransactionsRoot string   `json:"transactions_root"`
	StateRoot        string   `json:"state_root"`
	ReceiptsRoot     string   `json:"receipts_root"`
	Miner            string   `json:"miner"`
	Difficulty       *big.Int `json:"difficulty"`
	ExtraData        string   `json:"extra_data"`
	GasLimit         uint64   `json:"gas_limit"`
	GasUsed          uint64   `json:"gas_used"`
	Timestamp        uint64   `json:"timestamp"`
	Transactions     string   `json:"transactions"`
}

func NewClient(options client.Options) (client.Client, error) {
	if options.HostPort == "" {
		options.HostPort = os.Getenv("TEMPORAL_GRPC_ENDPOINT")
	}

	return client.NewClient(options)
}
