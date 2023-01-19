package tools

import (
	"context"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"log"
)

func GetWalletBalance(ctx context.Context, filename string, api lotusapi.FullNodeStruct) {
	walletlist := ReadFromConfig(filename)
	for _, k := range walletlist {
		add, _ := address.NewFromString(k)
		balance, _ := api.WalletBalance(ctx, add)
		log.Print(balance.String())

	}

}
