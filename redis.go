package main

import (
	"github.com/bingoohuang/go-utils"
	"github.com/go-redis/redis"
	"github.com/lunny/log"
	"strconv"
)

// password2/localhost:6388/0

type RedisServer struct {
	Addr      string
	Password  string
	DefaultDb int
}

func ParseServerItem(serverConfig string) *RedisServer {
	if serverConfig == "" {
		return nil
	}

	serverItems := go_utils.SplitTrim(serverConfig, "/")
	itemLen := len(serverItems)
	if itemLen == 1 {
		return &RedisServer{
			Addr:      serverItems[0],
			Password:  "",
			DefaultDb: 0,
		}
	} else if itemLen == 2 {
		dbIndex, _ := strconv.Atoi(serverItems[1])
		return &RedisServer{
			Addr:      serverItems[0],
			Password:  "",
			DefaultDb: dbIndex,
		}
	} else if itemLen == 3 {
		dbIndex, _ := strconv.Atoi(serverItems[2])
		return &RedisServer{
			Addr:      serverItems[1],
			Password:  serverItems[0],
			DefaultDb: dbIndex,
		}
	} else {
		log.Error("invalid servers argument")
		return nil
	}
}

func NewRedisClient(server *RedisServer) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     server.Addr,
		Password: server.Password,  // no password set
		DB:       server.DefaultDb, // use default DB
	})
}

func RedisGet(key string) (string, error) {
	client := NewRedisClient(redisServer)
	defer client.Close()

	return client.Get(key).Result()
}
