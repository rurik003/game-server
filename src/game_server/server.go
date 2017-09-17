package main

import (
	"bytes"
	"container/heap"
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
		select {
		case data, ok := <-c:
			if !ok {
				break
			}
			recv_str := string(data.txt[:])
			fmt.Println("Received: ", recv_str)
			s.dispatchMessage(data)
		case p := <-s.manager.roomChnl:
			key := p.addr
			_, ok := s.playerMap[key]
			if ok {
				delete(s.playerMap, key)
			}
		}

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
	key := m.addr
	_, ok := s.playerMap[key]
	if ok {
		s.notifyPlayers(m)
	} else {
		room := s.manager.queue.peek()
		if room.priority < RoomSize {
			fmt.Println("room check successful ", room.priority)
			s.sendMessageFromServer("room check successful", m)
			plr := makePlayer(room.priority, m.addr, room)
			room.plyrChnl <- plr

			s.playerMap[plr.addr] = plr
			room.priority++
			heap.Fix(&s.manager.queue, 0)
		} else {
			s.sendMessageFromServer("creating a new room", m)
			s.manager.createRoom(s, m)
		}
	}

}

func (s *GameServer) notifyPlayers(m Message) {

	room := s.playerMap[m.addr].room
	room.msgChnl <- m
}

func (s *GameServer) unknownMessage(m Message) {
	s.sendMessageFromServer("wtf did you send?", m)
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

func (s *GameServer) sendMessageFromServer(str string, m Message) {
	fmt.Println("writing message to player: ", str)

	var sendBuf bytes.Buffer

	curTime := time.Now().String()
	sendBuf.WriteString(curTime)
	sendBuf.WriteString(" ")
	sendBuf.WriteString(str)
	s.conn.WriteToUDP(sendBuf.Bytes(), m.addr)
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
