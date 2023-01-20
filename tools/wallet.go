package tools

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"log"
	"math/big"
	"net/http"
)

func GetWalletBalance(ctx context.Context, filename string, api lotusapi.FullNodeStruct) {
	walletlist := ReadFromConfig(filename)
	for _, k := range walletlist {
		add, _ := address.NewFromString(k)
		balance, _ := api.WalletBalance(ctx, add)
		balanceFIL := new(big.Int)
		balanceFIL.SetString(balance.String(), 10)
		balanceFIL.Div(balanceFIL, big.NewInt(1e18))

		// check if balance is less than 15
		if balanceFIL.Int64() >= 15.0 {
			log.Printf("钱包 %s 的余额为 %s FIL", add, balanceFIL)
		} else {
			//SendEm("钱包余额不足", []byte("钱包"+add.String()+"的余额为"+balanceFIL.Int64))
			SendEm(add.String(), []byte(add.String()+"的余额为"+balanceFIL.String()+"FIL"))
			message := fmt.Sprintf("钱包 %s 的余额为 %s FIL，不足 15 FIL", add, balanceFIL)
			resp, err := http.Post("https://oapi.dingtalk.com/robot/send?access_token=your_token", "application/json", bytes.NewBuffer([]byte(`{"msgtype": "text", "text": {"content": "`+message+`"}}`)))
			if err != nil {
				log.Printf("发送消息失败：%s", err)
			}
			defer resp.Body.Close()
		}

	}

}
