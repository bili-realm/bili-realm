package decode

import (
	"github.com/bytedance/sonic"
)

func Decode(data []byte) (*Payload, error) {
	var payload = new(Payload)
	node, err := sonic.Get(data, "cmd")
	if err != nil {
		return nil, err
	}
	if cmd, err := node.String(); err != nil {
		return nil, err
	} else {
		payload.Cmd = cmd
		switch cmd {
		case "DANMU_MSG":
			err := processDanmaku(payload, data)
			if err != nil {
				return nil, err
			}
		}
	}
	return payload, nil
}

func processDanmaku(payload *Payload, data []byte) (err error) {
	var info []interface{}
	if infoNode, err := sonic.Get(data, "info"); err == nil {
		if info, err = infoNode.Array(); err != nil {
			return err
		}
	} else {
		return err
	}
	// todo process info gracefully
	// todo user medal
	content := info[1].(string)
	timestamp := int64(info[0].([]any)[4].(float64))
	uid := int64(info[2].([]any)[0].(float64))
	uname := info[2].([]any)[1].(string)
	payload.Danmaku = &Danmaku{
		Content:   content,
		Timestamp: timestamp,
	}
	payload.User = &User{
		Uid:  uid,
		Name: uname,
	}
	return nil
}
