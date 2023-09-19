package tools

import (
	"context"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"log"
	"math/big"
)

var wallets = make(map[address.Address]*big.Int)
var notificationsSent = make(map[address.Address]bool)

func GetWalletBalance(ctx context.Context, filename string, api lotusapi.FullNodeStruct) {
	walletlist := ReadFromConfig(filename)
	for _, k := range walletlist {
		add, _ := address.NewFromString(k)
		balance, _ := api.WalletBalance(ctx, add)
		balanceFIL := new(big.Int)
		balanceFIL.SetString(balance.String(), 10)
		balanceFIL.Div(balanceFIL, big.NewInt(1e18))

		// 打印余额
		log.Printf("钱包 %s 的余额为 %s FIL", add, balanceFIL)

		// 检查余额是否低于15 FIL，并且未发送过通知
		if balanceFIL.Cmp(big.NewInt(15)) < 0 && !notificationsSent[add] {
			SendEm("余额不足", []byte(add.String()+"的余额为"+balanceFIL.String()+"FIL"))
			notificationsSent[add] = true
		} else if balanceFIL.Cmp(big.NewInt(15)) >= 0 {
			// 如果余额大于等于15 FIL，重置通知标志
			notificationsSent[add] = false
		}
	}
}
