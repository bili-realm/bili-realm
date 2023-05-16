package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Res struct {
	Code    int       `json:"code"`
	Data    DanmuInfo `json:"data"`
	Message string    `json:"message"`
	Ttl     int       `json:"ttl"`
}

type DanmuInfo struct {
	BusinessId       int        `json:"business_id"`
	Group            string     `json:"group"`
	HostList         []HostList `json:"host_list"`
	MaxDelay         int        `json:"max_delay"`
	RefreshRate      int        `json:"refresh_rate"`
	RefreshRowFactor float64    `json:"refresh_row_factor"`
	Token            string     `json:"token"`
}

type HostList struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	WsPort  int    `json:"ws_port"`
	WssPort int    `json:"wss_port"`
}

type VerifyPacket struct {
	Uid      int    `json:"uid"`
	RoomId   int    `json:"roomid"`
	ProtoVer int    `json:"protover"`
	Platform string `json:"platform"`
	Type     int    `json:"type"`
	Key      string `json:"key"`
}

const (
	_ = iota
	_
	HeartBeatPacketType
	HeartBeatResponsePacketType
	_
	NotifyPacketType
	_
	VerifyPacketType
	VerifyResponsePacketType
)

func getDanmukuInfo(roomId string) (string, []HostList, error) {
	r, err := http.Get("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=" + roomId)
	if err != nil {
		return "", nil, err
	}
	d := json.NewDecoder(r.Body)
	val := &Res{}
	err = d.Decode(&val)
	if err != nil {
		return "", nil, err
	}
	return val.Data.Token, val.Data.HostList, nil
}

func toPacket(payload []byte, packetType int) []byte {
	head := make([]byte, 16)
	binary.BigEndian.PutUint32(head[0:4], uint32(len(payload)+16))
	binary.BigEndian.PutUint16(head[4:6], uint16(16))
	binary.BigEndian.PutUint16(head[6:8], uint16(1))
	binary.BigEndian.PutUint32(head[8:12], uint32(packetType))
	binary.BigEndian.PutUint32(head[12:16], uint32(1))

	return append(head, payload...)
}

func main() {
	t, hl, err := getDanmukuInfo("7734200")
	if err != nil {
		log.Fatalf("req err: %v\n", err)
	}

	s := hl[0].Host
	i := hl[0].Port

	tcpClient, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s, strconv.Itoa(i)))
	if err != nil {
		log.Fatalf("tcp dial err: %#v\n", err)
	}

	verifyPacket := &VerifyPacket{
		Uid:      0,
		RoomId:   7734200,
		ProtoVer: 3,
		Platform: "danmuji",
		Type:     2,
		Key:      t,
	}

	b, err := json.Marshal(verifyPacket)
	if err != nil {
		log.Fatalf("json marshal err: %#v\n", err)
	}
	packet := toPacket(b, VerifyPacketType)

	fmt.Printf("packet: %s\nraw: %#v\n", packet, packet)

	_, err = tcpClient.Write(packet)
	if err != nil {
		log.Fatalf("tcp write err: %#v\n", err)
	}
	heartBeatPacket := toPacket([]byte("[object Object]"), HeartBeatPacketType)
	fmt.Printf("heartBearPacket: %#v\n", heartBeatPacket)
	_, err = tcpClient.Write(heartBeatPacket)
	if err != nil {
		log.Fatalf("tcp write err: %#v\n", err)
	}

	go func() {
		heartBeatTicker := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-heartBeatTicker.C:
				heartBeatPacket = toPacket([]byte("[object Object]"), HeartBeatPacketType)
				fmt.Printf("heartBearPacket: %s\n", heartBeatPacket)
				_, err := tcpClient.Write(heartBeatPacket)
				if err != nil {
					log.Fatalf("发送心跳失败: %v", err)
				}
				log.Print("发送心跳成功")
			}
		}
	}()

	var messageCh = make(chan []byte)

	var r []byte
	for {
		byteCount, err := tcpClient.Read(r)
		if err != nil {
			log.Fatalf("tcp read err: %#v\n", err)
		}
		if byteCount != 0 {
			log.Print(byteCount)
			messageCh <- r[:byteCount]
		}
	}

}
