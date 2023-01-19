package tools

import (
	"context"
	"github.com/filecoin-project/go-address"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"log"
	"strconv"
)

var f024972, f029401, f033123, f042540, f042558, f01785096, f01867066 uint64

func CheckPower(ctx context.Context, filename string, api lotusapi.FullNodeStruct, tipset types.TipSetKey) {

	minerlist := ReadFromConfig(filename)
	for _, k := range minerlist {
		maddr, _ := address.NewFromString(string(k))
		faults, _ := api.StateMinerFaults(context.Background(), maddr, tipset)
		count, _ := faults.Count()
		//fmt.Printf("Current chain head is: %s", tipset.String())
		//fmt.Print(faults.Count())
		log.Print(maddr.String(), "错误扇区数量为：", count)
		if count > 10 {

			switch maddr.String() {
			case "f024972":
				if f024972 < count {
					SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+strconv.FormatUint(count, 10)))
					f024972 = count
				}
			case "f029401":
				if f029401 < count {
					SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+strconv.FormatUint(count, 10)))
					f029401 = count
				}
			case "f033123":
				if f033123 < count {
					SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+strconv.FormatUint(count, 10)))
					f033123 = count
				}
			case "f042540":
				if f042540 < count {
					SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+strconv.FormatUint(count, 10)))
					f042540 = count
				}
			case "f042558":
				if f042558 < count {
					SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+strconv.FormatUint(count, 10)))
					f042558 = count
				}
			case "f01785096":
				if f01785096 < count {
					SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+strconv.FormatUint(count, 10)))
					f01785096 = count
				}
			case "f01867066":
				if f01867066 < count {
					SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+strconv.FormatUint(count, 10)))
					f01867066 = count
				}

			}
		}

	}
}
