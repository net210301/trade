package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/net210301/trade/pkg/tlog"
	"time"
)

type Redis struct {
	clientPool *redis.Pool
}

func newRedis(configure *Config) *Redis {
	return &Redis{&redis.Pool{
		MaxIdle:     configure.MaxIdle,
		MaxActive:   configure.MaxActive,
		IdleTimeout: time.Duration(configure.IdleTimeout) * time.Millisecond,
		Wait:        configure.Wait,
		Dial: func() (redis.Conn, error) {
			tlog.WatchDog.Debug("[CacheDatabase] Dial Connects To The CacheUtils Server")
			c, err := redis.Dial("tcp", configure.Address())
			if err != nil {
				return nil, err
			}
			if configure.Password != "" {
				if _, err := c.Do("AUTH", configure.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	},
	}
}

func (r *Redis) pingHealth() error {
	conn := r.GetClient()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return err
	}
	tlog.WatchDog.Info("[CacheDatabase] Connected to the Redis Database")
	return nil
}

func (r *Redis) GetClient() redis.Conn {
	return r.clientPool.Get()
}

func (r *Redis) SetString(key string, data string, time int) error {
	conn := r.clientPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}

	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Redis) Exists(key string) bool {
	conn := r.clientPool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func (r *Redis) Get(key string) ([]byte, error) {
	conn := r.clientPool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (r *Redis) Delete(key string) (bool, error) {
	conn := r.clientPool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func (r *Redis) LikeDeletes(key string) error {
	conn := r.clientPool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = r.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}