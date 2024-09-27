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

var sTime = time.Now()

func connectLotusAPI(ctx context.Context, addr, authToken string) (lotusapi.FullNodeStruct, jsonrpc.ClientCloser, error) {
	var api lotusapi.FullNodeStruct
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}

	var closer jsonrpc.ClientCloser
	var err error

	for {
		select {
		case <-ctx.Done():
			return api, nil, ctx.Err()
		default:
			closer, err = jsonrpc.NewMergeClient(ctx, "ws://"+addr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
			if err == nil {
				log.Printf("连接成功")
				return api, closer, nil
			}
			log.Printf("连接失败: %s, 5秒后重试...", err)
			time.Sleep(5 * time.Second)
			elapsedTime := time.Since(sTime)
			if elapsedTime > 10*time.Minute {
				log.Println("已经过了十分钟")
				tools.SendEm("连接失败", []byte(time.Now().String()))
				sTime = time.Now()
			}
		}
	}
}

func main() {
	home := os.Getenv("HOME")
	authToken := os.Getenv("TOKEN")
	addr := os.Getenv("ADDR")

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		api, closer, err := connectLotusAPI(ctx, addr, authToken)
		if err != nil {
			log.Printf("连接 Lotus API 失败: %v", err)
			cancel()
			time.Sleep(30 * time.Second)
			continue
		}

		log.Printf("Home 目录: %s", home)

		for {
			select {
			case <-ctx.Done():
				log.Println("上下文已取消，重新连接")
				closer()
				break
			default:
				tipset, err := api.ChainHead(ctx)
				if err != nil {
					log.Printf("获取链头失败: %v", err)
					time.Sleep(30 * time.Second)
					continue
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

				time.Sleep(30 * time.Second)
			}
		}

		closer()
		cancel()
	}
}
