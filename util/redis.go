package util

//redis连接操作
import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var Pool *redis.Pool

func InitPool(address string, maxIdle int, maxActive int, idleTimeout time.Duration) error {
	Pool = &redis.Pool{
		MaxIdle:     maxIdle,     //最大空闲连接数
		MaxActive:   maxActive,   //数据最大连接数,0没有限制
		IdleTimeout: idleTimeout, //最大空闲时间
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				address,
			)
		},
	}
	con := Pool.Get()
	if err := con.Err(); err != nil {
		fmt.Println("连接redis错误", err)
		return err
	}
	defer con.Close()
	return nil
}
