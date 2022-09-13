package Limiting_exercise

import (
	"context"
	"fmt"
	"time"
	"github.com/go-redis/redis/v9"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:       "localhost:6379",
		Password:   "",
		DB:         0,
		PoolSize:   3,
		MaxRetries: 3,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}


//@param key string object for rate limit such as uid/ip+url
//@param fillInterval time.Duration such as 1*time.Second
//@param limitNum max int64 allowed number per fillInterval
//@return whether reach rate limit, false means reach
func fixWindowRateLimit(key string, fillInterval time.Duration, limitNum int64) bool {

	ctx := context.Background()

	tick := int64(time.Now().Unix()) / int64(fillInterval.Seconds())

	currentKey := fmt.Sprintf("%s_%d_%d_%d", key, fillInterval, limitNum, tick)

	startCount := 0
	_, err := client.SetNx(ctx, currentKey, startCount, fillInterval).Result()
	if err != nil {
		panic(err)
	}

	quantum, err := client.Incr(ctx, currentKey).Result()
	if err != nil {
		panic(err)
	}

	if quantum > limitNum {
		return false
	}

	return true


}

//
func test1() {

	for i := 0; i < 10; i++ {
		go func() {
			rs := fixWindowRateLimit("test1", 1*time.Second, 5)
			fmt.Println("result is:", rs)
		}()
	}

}




