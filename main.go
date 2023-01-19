package post

import (
	"bufio"
	"context"
	"fmt"
	"github.com/filecoin-project/go-address"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
)

func main() {
	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.4tDmJiysQVzdMgpu70bvQHh1poD3pAv30MQsdW770fQ"
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}
	addr := "0.0.0.0:9999"

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

	ticker1 := time.NewTicker(5 * time.Second)
	// 一定要调用Stop()，回收资源
	defer ticker1.Stop()
	go func(t *time.Ticker) {
		for {
			// 每5秒中从chan t.C 中读取一次
			<-t.C
			fmt.Println("Ticker:", time.Now().Format("2006-01-02 15:04:05"))
			tipset, err := api.ChainHead(context.Background())
			if err != nil {
				log.Fatalf("calling chain head: %s", err)
			}
			br := bufio.NewReader(f)
			for {
				a, _, c := br.ReadLine()
				if c == io.EOF {
					break
				}
				maddr, _ := address.NewFromString(string(a))
				faults, _ := api.StateMinerFaults(context.Background(), maddr, tipset.Key())
				count, _ := faults.Count()
				//fmt.Printf("Current chain head is: %s", tipset.String())
				//fmt.Print(faults.Count())
				log.Print(maddr.String(), "错误扇区数量为：", count)
			}
		}
	}(ticker1)
	time.Sleep(30 * time.Second)
	fmt.Println("ok")

}
