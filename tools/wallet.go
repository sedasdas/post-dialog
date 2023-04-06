package tools

import (
	"context"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"log"
	"math/big"
)

var wallets = make(map[address.Address]*big.Int)

func GetWalletBalance(ctx context.Context, filename string, api lotusapi.FullNodeStruct) {
	walletlist := ReadFromConfig(filename)
	for _, k := range walletlist {
		add, _ := address.NewFromString(k)
		balance, _ := api.WalletBalance(ctx, add)
		balanceFIL := new(big.Int)
		balanceFIL.SetString(balance.String(), 10)
		balanceFIL.Div(balanceFIL, big.NewInt(1e18))
		wallets[add] = balanceFIL

		// check if balance is less than 15 FIL
		if balanceFIL.Cmp(big.NewInt(15)) >= 0 {
			log.Printf("钱包 %s 的余额为 %s FIL", add, balanceFIL)
		} else {
			if prevBalance, ok := wallets[add]; ok {
				if balanceFIL.Cmp(prevBalance) < 0 {
					wallets[add] = balanceFIL
					SendEm(add.String(), []byte(add.String()+"的余额为"+balanceFIL.String()+"FIL"))
					log.Printf("钱包 %s 的余额为 %s FIL，不足 15 FIL", add, balanceFIL)
				}
			} else {
				wallets[add] = balanceFIL
				SendEm(add.String(), []byte(add.String()+"的余额为"+balanceFIL.String()+"FIL"))
				log.Printf("钱包 %s 的余额为 %s FIL，不足 15 FIL", add, balanceFIL)
			}
		}
	}
}
