package tools

import (
	"context"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"log"
	"strconv"
	"sync"
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
		return err
	}
	count, err := faults.Count()
	if err != nil {
		return err
	}
	log.Printf("%s 错误扇区数量为：%d", miner.Address.String(), count)

	miner.FaultCount = count

	if miner.FaultCount != miner.LastAlertCount {
		if miner.FaultCount > 10 {
			if miner.FaultCount > miner.LastAlertCount {
				SendEm(miner.Address.String(), []byte(miner.Address.String()+"掉算力了，错误扇区数量为："+strconv.FormatUint(count, 10)))
			}
			if miner.FaultCount < miner.LastAlertCount {
				SendEm(miner.Address.String(), []byte(miner.Address.String()+"恢复中，错误扇区数量为："+strconv.FormatUint(count, 10)))
			}
		}
		miner.LastAlertCount = miner.FaultCount
	}
	return nil
}

func CheckPower(ctx context.Context, filename string, api lotusapi.FullNodeStruct, tipset types.TipSetKey) error {
	if miners == nil {
		if err := initMiners(filename); err != nil {
			return err
		}
	}

	var wg sync.WaitGroup
	for _, miner := range miners {
		wg.Add(1)
		go func(m *Miner) {
			defer wg.Done()
			if err := checkMinerPower(ctx, m, api, tipset); err != nil {
				log.Printf("检查矿工 %s 时出错: %v", m.Address.String(), err)
			}
		}(miner)
	}
	wg.Wait()
	return nil
}
