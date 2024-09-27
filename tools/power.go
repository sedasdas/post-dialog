package tools

import (
	"context"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"log"
)

type Miner struct {
	Address        address.Address
	FaultCount     uint64
	LastAlertCount uint64
}

var miners []*Miner

func initMiners(filename string) error {
	minerlist := ReadFromConfig(filename)
	miners = make([]*Miner, len(minerlist))
	for i, k := range minerlist {
		maddr, err := address.NewFromString(string(k))
		if err != nil {
			return err
		}
		miners[i] = &Miner{Address: maddr}
	}
	return nil
}

func checkMinerPower(ctx context.Context, miner *Miner, api lotusapi.FullNodeStruct, tipset types.TipSetKey) error {
	faults, err := api.StateMinerFaults(ctx, miner.Address, tipset)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	count, err := faults.Count()
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	log.Printf("%s 错误扇区数量为：%d", miner.Address.String(), count)

	miner.FaultCount = count
	return err
}

func CheckPower(ctx context.Context, filename string, api lotusapi.FullNodeStruct, tipset types.TipSetKey) error {
	if miners == nil {
		if err := initMiners(filename); err != nil {
			log.Printf(err.Error())
			return err
		}
	}
	for _, miner := range miners {

		if err := checkMinerPower(ctx, miner, api, tipset); err != nil {
			//log.Printf("检查矿工 %s 时出错: %v", miner.Address.String(), err)
		}
	}
	return nil
}
