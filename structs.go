package app

import (
	"math/big"
)

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

var TraceType = struct {
	Call         string
	Create       string
	DelegateCall string
	Reward       string
	Suicide      string
}{
	Call:         "call",
	Create:       "create",
	DelegateCall: "delegateCall",
	Reward:       "reward",
	Suicide:      "suicide",
}

type TraceBlockPayload struct {
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
}

type Trace struct {
	Action struct {
		CallType   string `json:"callType"`
		From       string `json:"from"`
		Gas        string `json:"gas"`
		Input      string `json:"input"`
		To         string `json:"to"`
		Value      string `json:"value"`
		RewardType string `json:"rewardType"`
	} `json:"action"`
	BlockHash   string `json:"blockHash"`
	BlockNumber uint64 `json:"blockNumber"`
	Result      struct {
		GasUsed string `json:"gasUsed"`
		Output  string `json:"output"`
	} `json:"result"`
	Subtraces           uint64        `json:"subtraces"`
	TraceAddress        []interface{} `json:"traceAddress"`
	TransactionHash     string        `json:"transactionHash"`
	TransactionPosition uint64        `json:"transactionPosition"`
	Type                string        `json:"type"`
	Error               string        `json:"error"`
}

type TraceBlockResponse struct {
	ID      int     `json:"id"`
	Jsonrpc string  `json:"jsonrpc"`
	Result  []Trace `json:"result"`
	Error   *struct {
		Code    int
		Message string
	}
}
