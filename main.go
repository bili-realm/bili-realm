package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/bili-realm/bili-realm/decode"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/url"
	"strconv"
	"time"
)

var danmakuUrl = url.URL{
	Host:   "broadcastlv.chat.bilibili.com",
	Scheme: "wss",
	Path:   "sub",
}

type Realm struct {
	RoomId       int
	Conn         *websocket.Conn
	RawPacket    chan *packet
	ParsedPacket chan *decode.Payload
	context.Context
}

func (r Realm) Start() context.CancelFunc {
	_, cancelFunc := context.WithCancel(r.Context)

	if verifyPacket, err := genVerifyPacket(verifyPacketBody{
		Uid:      0,
		RoomId:   r.RoomId,
		ProtoVer: 3,
		Platform: "web",
		Type:     2,
	}); err == nil {
		if err := r.Conn.WriteMessage(websocket.TextMessage, append(verifyPacket.head, verifyPacket.body...)); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}

	go func(conn *websocket.Conn) {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.TextMessage, genHeartbeatPacket().head); err != nil {
					panic(err)
				}
			}
		}
	}(r.Conn)

	go func() {
		for {
			_, msg, err := r.Conn.ReadMessage()
			if err != nil {
				panic(err)
			}
			handleMessage(r.RawPacket, msg)
		}
	}()

	return cancelFunc
}

func NewRealm(roomId int) *Realm {
	conn, _, err := websocket.DefaultDialer.Dial(danmakuUrl.String(), nil)
	if err != nil {
		panic(err)
	}
	return &Realm{
		RoomId:       roomId,
		Conn:         conn,
		RawPacket:    make(chan *packet),
		ParsedPacket: make(chan *decode.Payload),
		Context:      context.Background(),
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(fmt.Errorf("%v", err))
		}
	}()

	realm := NewRealm(7734200)
	shutdown := realm.Start()
	defer shutdown()
	i := 0
	for {
		select {
		case rawPacket := <-realm.RawPacket:
			fmt.Printf("packet %d: %v\n\n", i, string(rawPacket.body))
			i++
		}
	}
}

func handleMessage(rawCh chan<- *packet, msg []byte) {
	p := packet{
		head: msg[:16],
		body: msg[16:],
	}
	//packetLen := binary.BigEndian.Uint32(p.head[:4])
	//headLen := binary.BigEndian.Uint16(p.head[4:])
	packetVer := binary.BigEndian.Uint16(p.head[6:])
	packetTyp := binary.BigEndian.Uint32(p.head[8:])
	switch packetTyp {
	case verifyResponse:
		rawCh <- &p
		log.Println("进入房间")
	case heartbeatResponse:
		rawCh <- &packet{
			head: p.head,
			body: []byte(strconv.FormatUint(uint64(binary.BigEndian.Uint32(p.body)), 10)),
		}
		log.Println("人气值", binary.BigEndian.Uint32(p.body))
	case notify:
		if packetVer == 3 {
			brReader := brotli.NewReader(bytes.NewReader(p.body))
			if b, err := io.ReadAll(brReader); err == nil {
				packets := splitAndParse(b)
				for _, packet := range packets {
					rawCh <- &packet
				}
			} else {
				panic(err)
			}
		} else if packetVer == 0 {
			rawCh <- &p
		} else {
			panic(fmt.Sprintf("unknown packet version: %d\nraw bytes: %d\npacket body: %s\n", packetVer, p.body, string(p.body)))
		}
	default:
		panic(fmt.Sprintf("unknown packet type: %d\nraw bytes: %d\npacket body: %s\n", packetTyp, p.body, string(p.body)))
	}
}

func splitAndParse(data []byte) []packet {
	total := len(data)
	offset := 0
	var packets []packet
	for offset < total {
		l := binary.BigEndian.Uint32(data[offset:])
		packets = append(packets, packet{
			head: data[offset : offset+16],
			body: data[offset+16 : offset+int(l)],
		})
		offset += int(l)
	}
	return packets
}

func genVerifyPacket(p verifyPacketBody) (*packet, error) {
	if body, err := json.Marshal(p); err == nil {
		var head []byte = make([]byte, 16)
		binary.BigEndian.PutUint32(head, uint32(len(body)+16))
		binary.BigEndian.PutUint16(head[4:], 16)
		binary.BigEndian.PutUint16(head[6:], 1)
		binary.BigEndian.PutUint32(head[8:], verify)
		binary.BigEndian.PutUint32(head[12:], 1)

		return &packet{
			head: head,
			body: body,
		}, nil
	} else {
		return nil, err
	}
}

func genHeartbeatPacket() *packet {
	var head []byte = make([]byte, 16)
	binary.BigEndian.PutUint32(head, 16)
	binary.BigEndian.PutUint16(head[4:], 16)
	binary.BigEndian.PutUint16(head[6:], 1)
	binary.BigEndian.PutUint32(head[8:], heartbeat)
	binary.BigEndian.PutUint32(head[12:], 1)

	return &packet{
		head: head,
		body: []byte{},
	}
}
