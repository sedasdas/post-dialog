package tools

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"golang.org/x/sync/errgroup"
	_ "sync"
	"time"
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
			return fmt.Errorf("failed to parse miner address: %w", err)
		}
		miners[i] = &Miner{Address: maddr}
	}
	return nil
}

func retryOperation(ctx context.Context, operation func() error) error {
	backoff := time.Second
	for i := 0; i < 3; i++ {
		err := operation()
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			backoff *= 2
		}
	}
	return fmt.Errorf("operation failed after 3 retries")
}

func checkMinerPower(ctx context.Context, miner *Miner, api lotusapi.FullNodeStruct, tipset types.TipSetKey) error {
	//var faults types.BitField
	var count uint64

	err := retryOperation(ctx, func() error {
		var err error
		faults, err := api.StateMinerFaults(ctx, miner.Address, tipset)
		if err != nil {
			return fmt.Errorf("failed to get miner faults: %w", err)
		}
		count, err = faults.Count()
		if err != nil {
			return fmt.Errorf("failed to count faults: %w", err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("%s 错误扇区数量为：%d\n", miner.Address.String(), count)

	miner.FaultCount = count

	if miner.FaultCount != miner.LastAlertCount {
		if miner.FaultCount > 10 {
			message := fmt.Sprintf("%s", miner.Address.String())
			if miner.FaultCount > miner.LastAlertCount {
				message += fmt.Sprintf("掉算力了，错误扇区数量为：%d", count)
			} else {
				message += fmt.Sprintf("恢复中，错误扇区数量为：%d", count)
			}
			SendEm(miner.Address.String(), []byte(message))
		}
		miner.LastAlertCount = miner.FaultCount
	}
	return nil
}

func CheckPower(ctx context.Context, filename string, api lotusapi.FullNodeStruct, tipset types.TipSetKey) error {
	if miners == nil {
		if err := initMiners(filename); err != nil {
			return fmt.Errorf("failed to initialize miners: %w", err)
		}
	}

	g, ctx := errgroup.WithContext(ctx)
	sem := make(chan struct{}, 10) // Limit concurrent goroutines

	for _, miner := range miners {
		miner := miner // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return ctx.Err()
			}

			if err := checkMinerPower(ctx, miner, api, tipset); err != nil {
				return fmt.Errorf("failed to check miner %s power: %w", miner.Address.String(), err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error occurred while checking miner power: %w", err)
	}

	return nil
}
