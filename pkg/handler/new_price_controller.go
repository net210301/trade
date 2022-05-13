package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/net210301/trade/pkg/client/redis"
	"github.com/net210301/trade/pkg/pubsub"
	"github.com/net210301/trade/pkg/service"
	"net/http"
)

func GetLatestTrade(ctx *gin.Context){
	rc := redis.Client
	rk := fmt.Sprintf("streams=%s", "btcusdt@aggTrade")
	get, err := rc.Get(rk)
	if err != nil {
		return
	}
	var r BinanceData
	err = json.Unmarshal(get, &r)
	if err != nil{
		ctx.JSON(http.StatusOK,gin.H{
			"code": 10040,
			"error": "get price error",
		})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{
		"result": r,
	})
}

func GetLatestTradeNotification(ws service.WebSocketConn, err error) {
	borker := pubsub.DefaultBroker
	subscriber := borker.GetSubscriber("BTCUSD","BTCUSD")
	ch := subscriber.GetMessages()
	for {
		if msg, ok := <-ch; ok {
			err = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v", msg.GetPayload())))
		}
	}
}