package websocket

import (
	"github.com/togettoyou/wsc"
	"time"
)

type Websocket struct {
	Client 			*wsc.Wsc
	Done 			chan bool
	PingDuration 	time.Duration 	// 健康檢查間隔
}

func newWebSocketClient(config *Config) *Websocket {
	done := make(chan bool)
	return &Websocket{
		Client: wsc.New(config.DialUrl),
		Done:   done,
		PingDuration: config.PingDuration,
	}
}

