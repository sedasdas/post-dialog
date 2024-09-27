package tools

import (
	"context"
	"fmt"
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

// AsyncSendEm 异步发送消息
func AsyncSendEm(address string, message []byte) {
	go func() {
		SendEm(address, message)
	}()
}

func checkMinerPower(ctx context.Context, miner *Miner, api lotusapi.FullNodeStruct, tipset types.TipSetKey) error {
	faults, err := api.StateMinerFaults(ctx, miner.Address, tipset)
	if err != nil {
		log.Printf("Error getting miner faults: %v", err)
		return err
	}
	count, err := faults.Count()
	if err != nil {
		log.Printf("Error counting faults: %v", err)
		return err
	}
	log.Printf("%s 错误扇区数量为：%d", miner.Address.String(), count)

	miner.FaultCount = count

	if miner.FaultCount != miner.LastAlertCount {
		if miner.FaultCount > 10 {
			message := miner.Address.String()
			if miner.FaultCount > miner.LastAlertCount {
				message += "掉算力了，错误扇区数量为：" + strconv.FormatUint(count, 10)
			} else if miner.FaultCount < miner.LastAlertCount {
				message += "恢复中，错误扇区数量为：" + strconv.FormatUint(count, 10)
			}
			AsyncSendEm(miner.Address.String(), []byte(message))
		}
		miner.LastAlertCount = miner.FaultCount
	}

	return nil
}

func CheckPower(ctx context.Context, filename string, api lotusapi.FullNodeStruct, tipset types.TipSetKey) error {
	if miners == nil {
		if err := initMiners(filename); err != nil {
			log.Printf("Error initializing miners: %v", err)
			return err
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(miners))

	for _, miner := range miners {
		wg.Add(1)
		go func(m *Miner) {
			defer wg.Done()
			if err := checkMinerPower(ctx, m, api, tipset); err != nil {
				errChan <- fmt.Errorf("检查矿工 %s 时出错: %v", m.Address.String(), err)
			}
		}(miner)
	}

	wg.Wait()
	close(errChan)

	// 收集所有错误
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		// 这里你可以选择如何处理这些错误，例如记录日志或返回第一个错误
		for _, err := range errors {
			log.Printf("%v", err)
		}
		return fmt.Errorf("检查矿工时发生了 %d 个错误", len(errors))
	}

	return nil
}
