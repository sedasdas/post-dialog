package tools

import (
	"context"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"log"
	"strconv"
)

type Miner struct {
	Address    address.Address
	FaultCount uint64
	LastCount  uint64
}

func CheckPower(ctx context.Context, filename string, api lotusapi.FullNodeStruct, tipset types.TipSetKey) {
	minerlist := ReadFromConfig(filename)
	miners := make([]*Miner, len(minerlist))
	for i, k := range minerlist {
		maddr, _ := address.NewFromString(string(k))
		miners[i] = &Miner{Address: maddr}
	}
	for _, miner := range miners {
		faults, _ := api.StateMinerFaults(context.Background(), miner.Address, tipset)
		count, _ := faults.Count()
		log.Print(miner.Address.String(), "错误扇区数量为：", count)
		miner.FaultCount = count
		if miner.FaultCount > 10 && miner.FaultCount > miner.LastCount {
			SendEm(miner.Address.String(), []byte(miner.Address.String()+"错误扇区数量为："+strconv.FormatUint(count, 10)))
			miner.LastCount = miner.FaultCount
		}
	}
}
