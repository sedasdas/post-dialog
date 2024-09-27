package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"post-dialog/tools"
	"time"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
)

var sTime = time.Now()

func connectLotusAPI(ctx context.Context, addr, authToken string) (lotusapi.FullNodeStruct, jsonrpc.ClientCloser, error) {
	var api lotusapi.FullNodeStruct
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}

	closer, err := jsonrpc.NewMergeClient(ctx, "ws://"+addr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		return api, nil, fmt.Errorf("连接失败: %w", err)
	}

	log.Println("连接成功")
	return api, closer, nil
}

func performCheck(ctx context.Context, home string, api lotusapi.FullNodeStruct) {
	tipset, err := api.ChainHead(ctx)
	if err != nil {
		log.Printf("获取链头失败: %v", err)
		return
	}

	log.Printf("链高度: %d", tipset.Height())

	if err := tools.CheckPower(ctx, home+"/miner-list", api, tipset.Key()); err != nil {
		log.Printf("检查算力失败: %v", err)
	}

	// 如果需要，可以取消注释以下行
	// if err := tools.GetWalletBalance(ctx, home+"/wallet-list", api); err != nil {
	//     log.Printf("获取钱包余额失败: %v", err)
	// }
	// if err := tools.CheckNet(); err != nil {
	//     log.Printf("检查网络失败: %v", err)
	// }
}

func main() {
	home := os.Getenv("HOME")
	authToken := os.Getenv("TOKEN")
	addr := os.Getenv("ADDR")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		api, closer, err := connectLotusAPI(ctx, addr, authToken)
		if err != nil {
			log.Printf("连接 Lotus API 失败: %v", err)
			cancel()
			time.Sleep(30 * time.Second)
			continue
		}

		log.Printf("Home 目录: %s", home)

		func() {
			defer closer()
			defer cancel()

			for {
				select {
				case <-ticker.C:
					performCheck(ctx, home, api)
				case <-ctx.Done():
					log.Println("上下文已取消，重新连接")
					return
				}
			}
		}()

		// 如果这里退出了内部循环，说明需要重新连接
		time.Sleep(5 * time.Second)
	}
}
