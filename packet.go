package main

type packet struct {
	head []byte
	body []byte
}

type verifyPacketBody struct {
	Uid      int    `json:"uid"`
	RoomId   int    `json:"roomid"`
	ProtoVer int    `json:"protover"`
	Platform string `json:"platform"`
	Type     int    `json:"type"`
}

const (
	_ = iota
	_
	heartBeat
	heartBeatResponse
	_
	notify
	_
	verify
	verifyResponse
)
