package tools

import (
	"context"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"log"
	"math/big"
)

func GetWalletBalance(ctx context.Context, filename string, api lotusapi.FullNodeStruct) {
	walletlist := ReadFromConfig(filename)
	for _, k := range walletlist {
		add, _ := address.NewFromString(k)
		balance, _ := api.WalletBalance(ctx, add)
		balanceFIL := new(big.Int)
		balanceFIL.SetString(balance.String(), 10)
		balanceFIL.Div(balanceFIL, big.NewInt(1e18))
		log.Printf("钱包 %s 的余额为 %s FIL", add, balanceFIL)

	}

}
