# bili-realm

Well, this is a Danmakus tool for bili live, still under development.

## Usage

```bash
go get github.com/bili-realm/bili-realm
```

```go
package main

import (
	"fmt"
	brealm "github.com/bili-realm/bili-realm"
)

func main() {
	// create a new realm
	r := brealm.NewRealm(123456)
	// start the realm
	shutdown := r.Start()
	defer shutdown()
	go func() {
		// get the raw danmakus
		for danmaku := range r.RawPacket {
			fmt.Printf("%v", danmaku)
		}
	}()

	// get all the parsed message
	msgCh := r.ReceiveMessage()
	go func() {
		for msg := range msgCh {
			fmt.Println(msg)
		}
	}()
}

```