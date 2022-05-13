package websocket

import (
	"github.com/togettoyou/wsc"
	"time"
)

type Config struct {
	DialUrl 			string			// 連線網址
	PingDuration 		time.Duration 	// 健康檢查間隔
	WriteWait 			time.Duration	// 寫超時
	MaxMessageSize 		int64			// 消息最大得長度，預設 512 byte
	MinRecTime 			time.Duration	// 最小斷線崇廉時間間隔
	MaxRecTime			time.Duration	// 最大崇廉時間間隔
	RecFactor			float64			// 每次重連失敗後，繼續重新連線的乘數因子，退避指數
	MessageBufferSize	int				// Message 發送緩衝持大小
}


// DefaultConfig 正式專案會再加上讀取 config 的方法，這裡都只做預設的
func DefaultConfig()*Config{
	return &Config{
		DialUrl: "wss://stream.yshyqxx.com/stream?streams=btcusdt@aggTrade",
		PingDuration: 9 * time.Second,
		WriteWait: 10 * time.Second,
		MaxMessageSize: 2048,
		MinRecTime: 2 * time.Second,
		MaxRecTime: 60 * time.Second,
		RecFactor: 1.5,
		MessageBufferSize: 1024,
	}
}

// Build 給出 redis 的 instance
func (config *Config) Build() *Websocket {
	websocketClient := newWebSocketClient(config)
	websocketClient.Client.SetConfig(&wsc.Config{
		WriteWait: config.WriteWait,
		MaxMessageSize: config.MaxMessageSize,
		MinRecTime:config.MinRecTime,
		MaxRecTime: config.MaxRecTime,
		RecFactor: config.RecFactor,
		MessageBufferSize: config.MessageBufferSize,
	})
	return websocketClient
}