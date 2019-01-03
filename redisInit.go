package RedisInfluxdb

import (
	"encoding/json"
	"github.com/influxdata/platform/kit/errors"
	"gopkg.in/redis.v4"
	"math/rand"
	"strconv"
	"time"
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

func AddRedisData(addr string, password string, rediskey string, number int) error {
	cli, err := ConnRedis(addr, password)
	defer cli.Close()
	if err != nil {
		return err
	}

	_, err = cli.Del(rediskey).Result()
	if err != nil {
		return err
	}
	rand.Seed(time.Now().Unix())
	for i := 0; i < number; i++ {
		x := rand.Intn(10000)
		m := make(map[string]interface{})
		m["value"] = x
		j, _ := json.Marshal(m)
		cli.HSet(rediskey, "value"+strconv.Itoa(i), string(j))
	}
	return nil
}
