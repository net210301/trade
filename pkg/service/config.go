package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/net210301/trade/pkg/tlog"
	"go.uber.org/zap"
	"strings"
)

const ModName = "server.gin"

type Config struct {
	Host          string
	Port          int
	Mode          string
	logger *tlog.Logger
}

// DefaultConfig 正式專案會再加上讀取 config 的方法，這裡都只做預設的
func DefaultConfig()*Config{
	return &Config{
		Host: "0.0.0.0",
		Port:  8090,
		Mode: gin.ReleaseMode,
		logger: tlog.WatchDog.With(zap.String("mod name", strings.Replace(ModName, " ", ".", -1))),
	}
}

// Build 給出gin 的 instance 以及使用到的中間件，以及是否追蹤等等
func (config *Config) Build() *Server {
	server := newServer(config)
	server.Use(recoverMiddleware(config.logger))
	return server
}

func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}