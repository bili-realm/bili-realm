package main

import (
	"flag"
	"fmt"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "ws://localhost:8080/ws", "addr for http service")

func main() {
	flag.Parse()
	conn, r, err := websocket.DefaultDialer.Dial(*addr, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("r: %v\n", r)
	fmt.Printf("conn: %v\n", conn)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		fmt.Printf("p: %v\n", string(p))
	}
}
