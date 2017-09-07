package main

import (
	//	"bytes"
	//	"fmt"
	"net"
)

type GamePlayer struct {
	number     uint
	addr       *net.UDPAddr
	room       *GameRoom
	properties map[string]string
}

func makePlayer(n uint, a *net.UDPAddr, r *GameRoom) *GamePlayer {
	var retPlayer GamePlayer
	retPlayer.init(n, a, r)
	return &retPlayer
}

func (p *GamePlayer) init(n uint, a *net.UDPAddr, r *GameRoom) {
	p.number = n
	p.addr = a
	p.room = r
	p.properties = make(map[string]string)
}

func (p *GamePlayer) addProperty(key string, val string) {
	p.properties[key] = val
}
