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

func (p *GamePlayer) init(n uint, a *net.UDPAddr) {
	p.number = n
	p.addr = a
	p.properties = make(map[string]string)
}

func (p *GamePlayer) addProperty(key string, val string) {
	p.properties[key] = val
}
