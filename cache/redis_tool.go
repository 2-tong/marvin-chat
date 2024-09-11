package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"marvin-chat/config"
)

var rClient *redis.Client

func init() {
	rCfg := config.GetConfig().Redis
	rClient = redis.NewClient(&rCfg)
	ping := rClient.Ping(context.Background())
	pong, err := ping.Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong)
}
