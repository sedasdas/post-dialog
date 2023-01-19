package main

import (
	"context"
	"github.com/filecoin-project/go-address"
	"log"
	"net/http"
	"os"
	"post-dialog/tools"
	"time"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
)

func main() {
	var f024972, f029401, f033123, f042540, f042558, f01785096, f01867066 uint64
	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.4tDmJiysQVzdMgpu70bvQHh1poD3pAv30MQsdW770fQ"
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}
	addr := "10.0.1.93:9999"

	var api lotusapi.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+addr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()

	// Now you can call any API you're interested in.

	f, err := os.OpenFile("/home/lotus/miner-list", os.O_RDWR|os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	l := []string{"f024972", "f029401", "f033123", "f042540", "f042558", "f01785096", "f01867066"}
	for {
		time.Sleep(3 * time.Second)
		log.Print("我在定时执行任务")
		tipset, err := api.ChainHead(context.Background())
		if err != nil {
			log.Fatalf("calling chain head: %s", err)
		}

		log.Print(tipset.Height())
		for _, k := range l {
			maddr, _ := address.NewFromString(string(k))
			faults, _ := api.StateMinerFaults(context.Background(), maddr, tipset.Key())
			count, _ := faults.Count()
			//fmt.Printf("Current chain head is: %s", tipset.String())
			//fmt.Print(faults.Count())
			log.Print(maddr.String(), "错误扇区数量为：", count)
			if count > 100 {
				switch maddr.String() {
				case "f024972":
					if f024972 < count {
						tools.SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+string(count)))
						f024972 = count
					}
				case "f029401":
					if f029401 < count {
						tools.SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+string(count)))
						f029401 = count
					}
				case "f033123":
					if f033123 < count {
						tools.SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+string(count)))
						f033123 = count
					}
				case "f042540":
					if f042540 < count {
						tools.SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+string(count)))
						f042540 = count
					}
				case "f042558":
					if f042558 < count {
						tools.SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+string(count)))
						f042558 = count
					}
				case "f01785096":
					if f01785096 < count {
						tools.SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+string(count)))
						f01785096 = count
					}
				case "f01867066":
					if f01867066 < count {
						tools.SendEm(maddr.String(), []byte(maddr.String()+"错误扇区数量为："+string(count)))
						f01867066 = count
					}

				}
			}

		}
	}

}
