package main

//import (
//	"fmt"
//	"net"
//)

type GameRoom struct {
	players    []*GamePlayer
	chnl       chan Message
	properties map[string]string
}

func makeRoom(numPlayers uint) *GameRoom {
	var retRoom GameRoom
	retRoom.init(numPlayers)
	return &retRoom
}

func (r *GameRoom) init(numPlayers uint) {
	r.players = make([]*GamePlayer, numPlayers)
	r.properties = make(map[string]string)
	r.chnl = make(chan Message)
}

func (r *GameRoom) getPlayers() []*GamePlayer {
	return r.players
}

func (r *GameRoom) addPlayer(p *GamePlayer) {
	r.players = append(r.players, p)
	p.room = r
}
