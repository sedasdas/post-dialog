package tools

import (
	"context"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"log"
	"strconv"
)

func GetWalletBalance(ctx context.Context, filename string, api lotusapi.FullNodeStruct) {
	walletlist := ReadFromConfig(filename)
	for _, k := range walletlist {
		add, _ := address.NewFromString(k)
		balance, _ := api.WalletBalance(ctx, add)
		n, _ := strconv.ParseInt(balance.String(), 10, 64)
		// 将数字转换成字符串，并保留小数点后一位
		result := strconv.FormatFloat(float64(n)/1e18, 'f', 1, 64)
		log.Print(result)
	}

}
