package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"
)

func handleClient(conn *net.UDPConn) {
	var buf [1024]byte
	var sendBuf bytes.Buffer
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	curTime := time.Now().String()
	sendBuf.WriteString(curTime)
	sendBuf.WriteString(" hello from server")
	conn.WriteToUDP(sendBuf.Bytes(), addr)
}

func main() {
	fmt.Println("game-server")

	portString := ":1203"
	udpAddr, err := net.ResolveUDPAddr("udp4", portString)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	for {
		handleClient(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
