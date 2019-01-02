package RedisInfluxdb

import (
	"github.com/influxdata/platform/kit/errors"
	"gopkg.in/redis.v4"
)

func ConnRedis(addr string, password string) (*redis.Client, error) {
	if addr == "" {
		return nil, errors.New("redisaddr can not be empty")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	_, err := client.Ping().Result()
	//fmt.Println(pong, err)
	if err != nil {
		return nil, err
	}
	return client, nil
}
