package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/url"
)

var danmakuUrl = url.URL{
	Host:   "broadcastlv.chat.bilibili.com",
	Scheme: "wss",
	Path:   "sub",
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

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(fmt.Errorf("%v", err))
		}
	}()

	conn, _, err := websocket.DefaultDialer.Dial(danmakuUrl.String(), nil)
	if err != nil {
		panic(err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	if verifyPacket, err := genVerifyPacket(verifyPacketBody{
		Uid:      0,
		RoomId:   7734200,
		ProtoVer: 3,
		Platform: "web",
		Type:     2,
	}); err == nil {
		if err := conn.WriteMessage(websocket.TextMessage, append(verifyPacket.head, verifyPacket.body...)); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}

	// listen to danmaku
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		handleMessage(msg)
	}

}

func handleMessage(msg []byte) {
	p := packet{
		head: msg[:16],
		body: msg[16:],
	}
	//packetLen := binary.BigEndian.Uint32(p.head[:4])
	//headLen := binary.BigEndian.Uint16(p.head[4:])
	packetVer := binary.BigEndian.Uint16(p.head[6:])
	packetTyp := binary.BigEndian.Uint32(p.head[8:])
	switch packetTyp {
	case heartBeatResponse:
		println("heart beat response", binary.BigEndian.Uint32(p.body))
	case notify:
		if packetVer == 3 {
			brReader := brotli.NewReader(bytes.NewReader(p.body))
			if b, err := io.ReadAll(brReader); err == nil {
				//maybe multiple packets
				packets := splitAndParse(b)
				for _, packet := range packets {
					//todo decode packet
					println(string(packet.body))
				}
			} else {
				panic(err)
			}
		} else {
			panic(fmt.Sprintf("unknown packet version: %d", packetVer))
		}
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
