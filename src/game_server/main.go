package main

import (
	//	"bytes"
	"fmt"
	//	"net"
	//	"os"
	//	"time"
)

func main() {
	fmt.Println("game-server")
	server := GameServer{}
	server.init(":1203")
	server.handleRequests()
}
