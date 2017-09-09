package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type GameServer struct {
	conn      *net.UDPConn
	playerMap map[*net.UDPAddr]*GamePlayer
	mapMtx    sync.Mutex
	manager   RoomManager
}

type Message struct {
	txt  []byte
	addr *net.UDPAddr
}

func (s *GameServer) init(portString string) {
	udpAddr, err := net.ResolveUDPAddr("udp4", portString)
	checkError(err)
	s.conn, err = net.ListenUDP("udp", udpAddr)
	checkError(err)

	s.mapMtx = sync.Mutex{}
	s.playerMap = make(map[*net.UDPAddr]*GamePlayer)
	s.manager.init()
}

func (s *GameServer) handleRequests() {
	c := make(chan Message)

	go s.dispatch(c)
	go s.manager.manageRooms(s)
	firstRoom := s.manager.queue.peek()
	go s.runRoom(firstRoom)

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

		recv_str := string(data.txt[:])
		fmt.Println("Received: ", recv_str)
		s.dispatchMessage(data)
	}
}

func (s *GameServer) dispatchMessage(m Message) {
	switch m.txt[0] {
	case 'n':
		s.notifyPlayers(m)
	case 'c':
		s.createRoomAndPlayer(m)
	default:
		s.unknownMessage(m)
	}
}

func (s *GameServer) createRoomAndPlayer(m Message) {
	fmt.Println("I'm supposed to create a room")

	var sendBuf bytes.Buffer

	curTime := time.Now().String()
	sendBuf.WriteString(curTime)
	sendBuf.WriteString(" You just connected!")
	s.conn.WriteToUDP(sendBuf.Bytes(), m.addr)

	s.manager.dispChnl <- m
}

func (s *GameServer) notifyPlayers(m Message) {

	s.mapMtx.Lock()
	room := s.playerMap[m.addr].room
	s.mapMtx.Unlock()
	room.msgChnl <- m
}

func (s *GameServer) unknownMessage(m Message) {
	fmt.Println("Received Unknown Message: ", m.txt[0])

	var sendBuf bytes.Buffer

	curTime := time.Now().String()
	sendBuf.WriteString(curTime)
	sendBuf.WriteString(" wtf did you just do?")
	s.conn.WriteToUDP(sendBuf.Bytes(), m.addr)
}

func (s *GameServer) runRoom(r *GameRoom) {
	for {
		select {
		case msg := <-r.msgChnl:
			s.broadcast(msg.txt, r)
		case plyr := <-r.plyrChnl:
			fmt.Println("Received player with address ", plyr.addr)
			r.addPlayer(plyr)
		}
	}
}

func (s *GameServer) sendMessage(b []byte, plyr *GamePlayer) {
	s.conn.WriteToUDP(b, plyr.addr)
}

func (s *GameServer) broadcast(b []byte, r *GameRoom) {
	fmt.Println("broadcast called")
	for i := range r.players {
		s.sendMessage(b, r.players[i])
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
