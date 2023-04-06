package main

import (
	"context"
	"log"
	"net/http"
	"post-dialog/tools"
	"time"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
)

func main() {

	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.4tDmJiysQVzdMgpu70bvQHh1poD3pAv30MQsdW770fQ"
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}
	addr := "127.0.0.1:9999"

	var api lotusapi.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+addr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()

	// Now you can call any API you're interested in.

	//l := []string{"f024972", "f029401", "f033123", "f042540", "f042558", "f01785096", "f01867066"}
	for {
		time.Sleep(30 * time.Second)
		tipset, err := api.ChainHead(context.Background())
		if err != nil {
			log.Fatalf("calling chain head: %s", err)
		}
		log.Print(tipset.Height())
		tools.CheckPower(context.Background(), "/home/lotus/miner-list", api, tipset.Key())
		tools.GetWalletBalance(context.Background(), "/home/lotus/wallet-list", api)
		tools.CheckNet()

	}

}
