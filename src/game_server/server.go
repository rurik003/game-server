package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"
)

type GameServer struct {
	conn *net.UDPConn
}

func (s *GameServer) init(portString string) {
	udpAddr, err := net.ResolveUDPAddr("udp4", portString)
	checkError(err)
	s.conn, err = net.ListenUDP("udp", udpAddr)
	checkError(err)
}

func (s *GameServer) handleRequests() {
	for {
		s.handleClient()
	}
}

func (s *GameServer) handleClient() {
	var buf [1024]byte
	var sendBuf bytes.Buffer
	_, addr, err := s.conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	curTime := time.Now().String()
	sendBuf.WriteString(curTime)
	sendBuf.WriteString(" hello from server")
	s.conn.WriteToUDP(sendBuf.Bytes(), addr)
}

func (s *GameServer) sendMessage(b []byte, other GamePlayer) {
	s.conn.WriteToUDP(b, other.addr)
}

func (s *GameServer) broadcast(sender uint, b []byte, r GameRoom) {
	players := r.getPlayers()
	for i := range players {
		if players[i].number != sender {
			s.sendMessage(b, players[i])
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
