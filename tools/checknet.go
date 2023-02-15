package tools

import (
	"log"
	"net"
)

func CheckNet() {
	_, err := net.Dial("tcp", "43.228.180.125:3453")
	if err != nil {
		log.Println("网络中断")
		SendEm("网络中断", []byte("网络中断"))
	}
}
