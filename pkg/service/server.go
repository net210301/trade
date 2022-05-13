package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type Server struct {
	*gin.Engine
	MyServer *http.Server
	config *Config
	listener net.Listener
	startTime time.Time
}

// Stop 直接終止Server
func (s *Server) Stop() error {
	return s.MyServer.Close()
}

// GracefulStop 不接受新的 request，把舊有的消化完之後再關閉
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.MyServer.Shutdown(ctx)
}

//Upgrade protocol to WebSocket
func (s *Server) Upgrade(ws *WebSocket) gin.IRoutes {
	return s.GET(ws.Pattern, func(c *gin.Context) {
		ws.Upgrade(c.Writer, c.Request)
	})
}

// Serve 服務啟動
func (s *Server) Serve() error {
	serveTime := time.Now()
	for _, route := range s.Engine.Routes() {
		s.config.logger.Info("[WebServer] Register Route", zap.String("method",route.Method), zap.String("path",route.Path))
	}
	s.MyServer = &http.Server{
		Addr:    s.config.Address(),
		Handler: s,
	}
	var err error
	s.config.logger.Info(fmt.Sprintf("[WebServer] Serving HTTP(%s) in %dms",s.config.Address() ,serveTime.Sub(s.startTime).Milliseconds()))
	err = s.MyServer.Serve(s.listener)

	if err == http.ErrServerClosed {
		s.config.logger.Warn("[WebServer] Close Gin :", zap.String("address",s.config.Address()))
		return nil
	}
	return err
}

func newServer(configure *Config) *Server{
	listener, err := net.Listen("tcp", configure.Address())
	if err != nil {
		configure.logger.Panic("Init gin server err", zap.Error(err))
	}
	configure.Port = listener.Addr().(*net.TCPAddr).Port
	gin.SetMode(configure.Mode)
	return &Server{
		Engine:   gin.New(),
		config:   configure,
		listener: listener,
		startTime: time.Now(),
	}
}