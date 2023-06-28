package decode

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
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

func processDanmaku(payload *Payload, data []byte) error {
	var err error

	var info []interface{}
	var infoNode ast.Node
	if infoNode, err = sonic.Get(data, "info"); err == nil {
		if info, err = infoNode.Array(); err != nil {
			return err
		}
	} else {
		return err
	}
	// todo process info gracefully
	content := info[1].(string)
	timestamp := int64(info[0].([]any)[4].(float64))
	uid := int64(info[2].([]any)[0].(float64))
	uname := info[2].([]any)[1].(string)
	// todo user medal
	payload.Payload = Danmaku{
		Content: content,
		User: User{
			Uid:  uid,
			Name: uname,
		},
		Timestamp: timestamp,
	}
	fmt.Printf("process danmaku: %#v\n", payload)
	return nil
}
