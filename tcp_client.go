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

func main() {
	t, hl, err := getDanmukuInfo("7734200")
	if err != nil {
		log.Fatalf("req err: %v\n", err)
	}

	s := hl[0].Host
	i := hl[0].Port

	// raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", s, strconv.Itoa(i)))
	// if err != nil {
	// 	log.Fatalf("ResolveTcp Err: %#v", err)
	// }
	// fmt.Printf("raddr: %v\n", raddr)
	// tcpClient, err = net.DialTCP("tcp", nil, raddr)
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

	head := make([]byte, 16)

	b, err := json.Marshal(verifyPacket)
	binary.BigEndian.PutUint32(head[0:4], uint32(len(b)+16))
	binary.BigEndian.PutUint16(head[4:6], uint16(16))
	binary.BigEndian.PutUint16(head[6:8], uint16(1))
	binary.BigEndian.PutUint32(head[8:12], uint32(7))
	binary.BigEndian.PutUint32(head[12:16], uint32(1))

	packet := append(head, b...)
	fmt.Printf("%v\n", head)
	fmt.Printf("packet: %s\nlen: %d", packet, len(packet))

	_, err = tcpClient.Write(packet)
	if err != nil {
		log.Fatalf("tcp write err: %#v\n", err)
	}
	go func() {
		var r []byte
		for {
			byteCount, err := tcpClient.Read(r)
			if err != nil {
				log.Fatalf("tcp read err: %#v\n", err)
			}
			if byteCount != 0 {
				fmt.Printf("byteCount: %v\n", byteCount)
				fmt.Printf("string(r): %v\n", string(r))
			}
		}
	}()

	heartBeatTicker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-heartBeatTicker.C:
			heartBeatPacket := make([]byte, 16)
			_, err := tcpClient.Write(heartBeatPacket)
			if err != nil {
				log.Fatalf("发送心跳失败: %v", err)
			}
			log.Print("发送心跳成功")
		}
	}
}
