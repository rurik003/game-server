package main

//import (
//	"fmt"
//	"net"
//)

type GameRoom struct {
	players    []GamePlayer
	chnl       chan Message
	properties map[string]string
}

func (r *GameRoom) init(numPlayers uint) {
	r.players = make([]GamePlayer, numPlayers)
	r.properties = make(map[string]string)
	r.chnl = make(chan Message)
}

func (r *GameRoom) getPlayers() []GamePlayer {
	return r.players
}
