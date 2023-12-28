package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"post-dialog/tools"
	"time"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
)

func connectLotusAPI(addr, authToken string) (lotusapi.FullNodeStruct, jsonrpc.ClientCloser, error) {
	var api lotusapi.FullNodeStruct
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}

	var closer jsonrpc.ClientCloser
	var err error

	for {
		closer, err = jsonrpc.NewMergeClient(context.Background(), "ws://"+addr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
		if err == nil {
			log.Printf("连接成功")
			break
		}
		log.Printf("连接失败: %s, retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	return api, closer, err
}

func main() {
	home := os.Getenv("HOME")
	authToken := os.Getenv("TOKEN")
	addr := os.Getenv("ADDR")

	api, closer, _ := connectLotusAPI(addr, authToken)
	defer closer()

	for {
		tipset, err := api.ChainHead(context.Background())
		if err != nil {
			log.Printf("发生故障: %s", err)
			closer()
			api, closer, err = connectLotusAPI(addr, authToken)
			continue
		}

		log.Printf("chain head: %d", tipset.Height())
		tools.CheckPower(context.Background(), home+"/miner-list", api, tipset.Key())
		tools.GetWalletBalance(context.Background(), home+"/wallet-list", api)
		tools.CheckNet()

		time.Sleep(30 * time.Second)
	}
}
