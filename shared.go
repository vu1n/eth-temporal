package app

import (
	"math/big"
)

const NewBlockTaskQueue = "NEW_BLOCK_TASK_QUEUE"
const BackfillTaskQueue = "BACKFILL_TASK_QUEUE"

var QueryTypes = struct {
	GET_NEXT_BLOCK string
	BLOCK_COMPLETE string
}{
	GET_NEXT_BLOCK: "get_next_block",
	BLOCK_COMPLETE: "block_complete",
}

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
	Number           uint64
	Hash             string
	ParentHash       string
	Sha3Uncles       string
	TransactionsRoot string
	StateRoot        string
	ReceiptsRoot     string
	Miner            string
	Difficulty       *big.Int
	ExtraData        string
	GasLimit         uint64
	GasUsed          uint64
	Timestamp        uint64
	Transactions     string
}
