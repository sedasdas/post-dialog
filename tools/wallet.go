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
		mfilFloat, _ := new(big.Float).SetString(balance.String())
		filFloat := new(big.Float).Mul(mfilFloat, big.NewFloat(1000))

		filString, _ := filFloat.Float64()
		log.Printf("钱包 %s 的余额为 %s", add, filString)

	}

}
