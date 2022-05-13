package redis

import (
	"fmt"
)

type Config struct {
	Host        string `default:"127.0.0.1" required:"true"`
	Port        int    `default:"6379" required:"true"`
	Password    string `default:""`
	MaxIdle     int    `default:"30"`
	MaxActive   int    `default:"20"`
	IdleTimeout int    `default:"200"`
	Wait        bool   `default:"true"`

}

// DefaultConfig 正式專案會再加上讀取 config 的方法，這裡都只做預設的
func DefaultConfig()*Config{
	return &Config{
		Host: "127.0.0.1",
		Port: 6379,
		MaxIdle: 30,
		MaxActive: 20,
		IdleTimeout: 200,
		Wait: true,
	}
}

// Build 給出 redis 的 instance
func (config *Config) Build() *Redis {
	redisClientPool := newRedis(config)
	return redisClientPool
}

func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}