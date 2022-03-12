package main

import (
	"bytes"
	"encoding/json"
	"eth-temporal/app"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PayloadTraceBlock struct {
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
}

func main() {
	postBody, _ := json.Marshal(PayloadTraceBlock{
		Method:  "trace_block",
		Params:  []string{fmt.Sprintf("0x%x", 14357125)},
		Id:      1,
		Jsonrpc: "2.0",
	})
	// fmt.Print(string(postBody))
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
	fmt.Println("Length: ", len(res.Result))
	fmt.Println(res.Result[0].BlockHash)
	// _ = os.WriteFile("test.json", body, 0644)
	// sb := string(body)
	// fmt.Print(sb)
}
