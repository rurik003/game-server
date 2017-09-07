package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"
)

type GameServer struct {
	conn  *net.UDPConn
	rooms []GameRoom
}

type Message struct {
	msg  []byte
	addr *net.UDPAddr
}

func (s *GameServer) init(portString string) {
	udpAddr, err := net.ResolveUDPAddr("udp4", portString)
	checkError(err)
	s.conn, err = net.ListenUDP("udp", udpAddr)
	checkError(err)
}

func (s *GameServer) handleRequests() {
	c := make(chan Message)

	go s.dispatch(c)
	s.handleClient(c)
}

func (s *GameServer) handleClient(c chan Message) {
	for {
		var buf [1024]byte

		sz, addr, err := s.conn.ReadFromUDP(buf[0:])
		if err != nil {
			return
		}
		msg := Message{buf[:sz], addr}
		c <- msg
	}
}

func (s *GameServer) dispatch(c chan Message) {
	for {
		data, ok := <-c

		if !ok {
			return
		}

		recv_str := string(data.msg[:])
		fmt.Println("Received: ", recv_str)
		s.dispatchMessage(data)
	}
}

func (s *GameServer) dispatchMessage(m Message) {
	switch m.msg[0] {
	case 'n':
		s.notifyPlayers(m)
	case 'c':
		s.createRoomAndPlayer(m)
	default:
		s.unknownMessage(m)
	}
}

func (s *GameServer) notifyPlayers(m Message) {
	fmt.Println("I'm supposed to be notifying the players")

	var sendBuf bytes.Buffer

	curTime := time.Now().String()
	sendBuf.WriteString(curTime)
	sendBuf.WriteString(" You just tried to notify players!")
	s.conn.WriteToUDP(sendBuf.Bytes(), m.addr)
}

func (s *GameServer) createRoomAndPlayer(m Message) {
	fmt.Println("I'm supposed to create a room")

	var sendBuf bytes.Buffer

	curTime := time.Now().String()
	sendBuf.WriteString(curTime)
	sendBuf.WriteString(" You just tried to connect!")
	s.conn.WriteToUDP(sendBuf.Bytes(), m.addr)
}

func (s *GameServer) unknownMessage(m Message) {
	fmt.Println("Received Unknown Message: ", m.msg[0])

	var sendBuf bytes.Buffer

	curTime := time.Now().String()
	sendBuf.WriteString(curTime)
	sendBuf.WriteString(" wtf did you just do?")
	s.conn.WriteToUDP(sendBuf.Bytes(), m.addr)
}

func (s *GameServer) sendMessage(b []byte, other GamePlayer) {
	s.conn.WriteToUDP(b, other.addr)
}

func (s *GameServer) broadcast(sender uint, b []byte, r GameRoom) {
	players := r.getPlayers()
	for i := range players {
		s.sendMessage(b, players[i])
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
