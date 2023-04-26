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

func main() {
	home := os.Getenv("HOME")
	authToken := os.Getenv("TOKEN")
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}
	addr := os.Getenv("ADDR")

	//var minerapi lotusapi.StorageMiner

	var api lotusapi.FullNodeStruct

	// 连接 Lotus API
	var closer jsonrpc.ClientCloser
	var err error
	for {
		closer, err = jsonrpc.NewMergeClient(context.Background(), "ws://"+addr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
		if err == nil {
			log.Printf("connected to lotus successfully")
			break
		}
		log.Printf("connecting with lotus failed: %s, retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	defer closer()

	// 使用 Lotus API
	for {
		tipset, err := api.ChainHead(context.Background())
		if err != nil {
			log.Printf("calling chain head: %s", err)
			closer()
			for {
				closer, err = jsonrpc.NewMergeClient(context.Background(), "ws://"+addr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
				if err == nil {
					log.Printf("reconnected to lotus successfully")

					break
				}
				log.Printf("reconnecting with lotus failed: %s, retrying in 5 seconds...", err)
				time.Sleep(5 * time.Second)
			}
			continue
		}
		log.Printf("chain head: %d", tipset.Height())
		tools.CheckPower(context.Background(), home+"/miner-list", api, tipset.Key())
		tools.GetWalletBalance(context.Background(), home+"/wallet-list", api)
		tools.CheckNet()
		// 在这里使用 Lotus API
		time.Sleep(30 * time.Second)
	}
}
