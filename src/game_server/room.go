package main

import (
	"container/heap"
)

const RoomSize = 8
const NumRooms = 8

type GameRoom struct {
	priority   uint
	players    []*GamePlayer
	msgChnl    chan Message
	plyrChnl   chan *GamePlayer
	properties map[string]string
}

//priority queue for matchmaking
type RoomQueue []*GameRoom

type RoomManager struct {
	dispChnl chan Message
	roomChnl chan *GamePlayer
	queue    RoomQueue
}

func (r *RoomManager) init() {
	r.dispChnl = make(chan Message)
	r.roomChnl = make(chan *GamePlayer)

	//start with 1 room in priority queue
	r.queue = make(RoomQueue, 1, NumRooms)
	newRoom := makeRoom(RoomSize)
	r.queue[0] = newRoom
	heap.Init(&r.queue)
}

func (r *RoomManager) createRoom(s *GameServer, m Message) {
	newRoom := makeRoom(RoomSize)
	newPlayer := makePlayer(0, m.addr, newRoom)

	newRoom.addPlayer(newPlayer)
	newRoom.priority++
	heap.Push(&r.queue, newRoom)
	go s.runRoom(newRoom)

	s.playerMap[newPlayer.addr] = newPlayer
}

func makeRoom(numPlayers uint) *GameRoom {
	var retRoom GameRoom
	retRoom.init(numPlayers)
	return &retRoom
}

func (r *GameRoom) init(numPlayers uint) {
	r.priority = 0
	r.players = make([]*GamePlayer, 0, numPlayers)
	r.properties = make(map[string]string)
	r.msgChnl = make(chan Message)
	r.plyrChnl = make(chan *GamePlayer)
}

func (r *GameRoom) addPlayer(p *GamePlayer) {
	r.players = append(r.players, p)
	p.room = r
}

func (rq RoomQueue) Len() int {
	return len(rq)
}

func (rq RoomQueue) Less(i, j int) bool {
	return rq[i].priority < rq[j].priority
}

func (rq RoomQueue) Swap(i, j int) {
	rq[i], rq[j] = rq[j], rq[i]
}

func (rq *RoomQueue) Push(r interface{}) {
	room := r.(*GameRoom)
	*rq = append(*rq, room)
}

func (rq *RoomQueue) Pop() interface{} {
	old := *rq
	n := len(old)
	room := old[n-1]
	*rq = old[0 : n-1]
	return room
}

func (rq RoomQueue) peek() *GameRoom {
	return rq[0]
}
