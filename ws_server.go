package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ws(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("r.Header: %v\n", r.Header)
	c, err := upgrader.Upgrade(w, r, nil)
	t := time.NewTicker(time.Second)
	if err != nil {
		log.Fatalf("err: %v\n", err)
		return
	}
	defer c.Close()
	for {
		select {
		case t := <-t.C:
			err = c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Fatalf("err: %v\n", err)
			}
		}
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/ws", ws)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
