package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

func NewClient(options client.Options) (client.Client, error) {
	if options.HostPort == "" {
		options.HostPort = os.Getenv("TEMPORAL_GRPC_ENDPOINT")
	}

	return client.NewClient(options)
}

func HexToUInt(hexStr string) uint64 {
	trimmed := strings.Replace(hexStr, "0x", "", -1)
	trimmed = strings.Replace(trimmed, "0X", "", -1)
	result, _ := strconv.ParseUint(trimmed, 16, 64)
	return result
}

func HexToFloat(hexStr string) float32 {
	return float32(HexToUInt(hexStr))
}

func UIntToHex(value uint64) string {
	return fmt.Sprintf("0x%x", value)
}
