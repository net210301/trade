package handler

import (
	"encoding/json"
	"fmt"
	gw "github.com/gorilla/websocket"
	"github.com/net210301/trade/pkg/client/redis"
	"github.com/net210301/trade/pkg/client/websocket"
	"github.com/net210301/trade/pkg/pubsub"
	"github.com/net210301/trade/pkg/tlog"
	"sync"
	"time"
)

type Result struct {
	code int 		`json:"code"`
	msg  string 	`json:"message"`
}

type BinanceReceive struct {
	Stream 	string 	`json:"stream"`
	Data 	BinanceData 	`json:"data"`
}

type BinanceData struct {
	Ee	string 	`json:"e"`
	E	int64 	`json:"E"`
	S	string	`json:"s"`
	A	int64	`json:"a"`
	P	string	`json:"p"`
	Q	string	`json:"q"`
	F 	int64	`json:"f"`
	L 	int64	`json:"l"`
	T 	int64	`json:"T"`
	Mm 	bool	`json:"m"`
	M 	bool	`json:"M"`
}

func ProcessStreamBTCUSDT() {
	binanceWebsocket := websocket.BinanceWsClient
	// 處理連線, 每 N 秒發一次健康檢查連線
	binanceWebsocket.Client.OnConnected(func() {
		var wg sync.WaitGroup
		tlog.WatchDog.Info(fmt.Sprintf("OnConnected: %s",binanceWebsocket.Client.WebSocket.Url))
		wg.Add(1)
		go func() {
			defer wg.Done()
			t := time.NewTicker(binanceWebsocket.PingDuration)
			for {
				select {
				case <-t.C:
					_ = binanceWebsocket.Client.WebSocket.Conn.WriteMessage(gw.PingMessage, []byte(""))
				}
			}
		}()
	})

	// 錯誤處理
	binanceWebsocket.Client.OnDisconnected(func(err error) {
		tlog.WatchDog.Error(fmt.Sprintf("OnDisconnected: %s",err.Error()))
	})

	// 關閉前要做的事
	binanceWebsocket.Client.OnClose(func(code int, text string) {
		closeMSg := &Result{
			code: code,
			msg: text,
		}
		tlog.WatchDog.Error(fmt.Sprintf("OnClose: %v",closeMSg))
		binanceWebsocket.Done <- true
	})

	binanceWebsocket.Client.OnPingReceived(func(appData string) {
		tlog.WatchDog.Info("OnPingReceived")
	})

	binanceWebsocket.Client.OnPongReceived(func(appData string) {
		tlog.WatchDog.Info(fmt.Sprintf("OnPongReceived: %s",appData))
	})

	binanceWebsocket.Client.OnTextMessageReceived(func(message string) {
		// Save to Redis
		var r BinanceReceive
		err := json.Unmarshal([]byte(message), &r)
		if err != nil {
			tlog.WatchDog.Info(fmt.Sprintf("error: %s", err))
		}
		rc := redis.Client
		rk := fmt.Sprintf("streams=%s", r.Stream)

		b, err := json.Marshal(r.Data)
		if err != nil {
			tlog.WatchDog.Info(fmt.Sprintf("error: %s", err))
			return
		}

		err = rc.SetString(rk, string(b), -1)
		if err != nil {
			return
		}
		pubsub.DefaultBroker.Broadcast(string(b),"BTCUSD")
	})
	go binanceWebsocket.Client.Connect()
	for {
		select {
		case <-binanceWebsocket.Done:
			return
		}
	}
}