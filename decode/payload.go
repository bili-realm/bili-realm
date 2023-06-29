package decode

import "time"

type Payload struct {
	Cmd string
	*User
	*Medal
	*Danmaku
	*SuperChat
	*Gift
}

type User struct {
	Uid   int64
	Name  string
	Medal Medal
}

type Medal struct {
	Level  int
	Name   string
	Host   string
	HostId int
}

type Danmaku struct {
	Extra     string
	Emoticon  string
	Content   string
	Timestamp int64
}

type SuperChat struct {
	Danmaku
	Price    float64
	KeepTime time.Duration
}

type Gift struct {
	Name  string
	Count int
	User  User
	Price int
}
