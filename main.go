package main

import (
	"github.com/net210301/trade/pkg/handler"
	"github.com/net210301/trade/pkg/service"
)


func main(){
	httpCfg := service.DefaultConfig()
	server := httpCfg.Build()
	go handler.ProcessStreamBTCUSDT()
	server.GET("/get_latest_trade", handler.GetLatestTrade)
	server.Upgrade(service.WebSocketOptions("/get_latest_trade_notification",handler.GetLatestTradeNotification))
	server.Serve()
}
