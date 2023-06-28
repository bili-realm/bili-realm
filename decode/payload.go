package decode

type Payload struct {
	Cmd     string
	Payload any
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
	User      User
	Timestamp int64
}

type Gift struct {
	Name  string
	Count int
	User  User
	Price int
}
