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
		tools.SendEm("连接失败", []byte(time.Now().String()))

	}

	return api, closer, err
}

func main() {
	home := os.Getenv("HOME")
	authToken := os.Getenv("TOKEN")
	addr := os.Getenv("ADDR")

	api, closer, _ := connectLotusAPI(addr, authToken)
	defer closer()

	// 用于控制发送邮件的时间间隔
	mailTicker := time.NewTicker(10 * time.Minute)
	defer mailTicker.Stop()

	// 初始设定为 true，表示刚启动程序时可以发送邮件
	canSendMail := true

	for {
		tipset, err := api.ChainHead(context.Background())
		if err != nil {
			log.Printf("发生故障: %s", err)
			closer()

			// 如果可以发送邮件（即满足条件）
			if canSendMail {
				tools.SendEm("连接失败", []byte(time.Now().String()))
				canSendMail = false // 设置为 false，表示不可以发送邮件
			}

			api, closer, err = connectLotusAPI(addr, authToken)
			continue
		}

		log.Printf("chain head: %d", tipset.Height())
		tools.CheckPower(context.Background(), home+"/miner-list", api, tipset.Key())
		tools.GetWalletBalance(context.Background(), home+"/wallet-list", api)
		tools.CheckNet()

		select {
		case <-mailTicker.C:
			// 到达间隔时间后，设置可以发送邮件的条件为 true，以便下一次故障时可以再次发送邮件
			canSendMail = true
		default:
			// 如果还没到间隔时间，继续循环
		}

		time.Sleep(30 * time.Second)
	}
}
