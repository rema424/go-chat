package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool = &redis.Pool{
	MaxIdle:     3,
	IdleTimeout: 240 * time.Second,
	Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", ":6379") },
	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		if time.Since(t) < time.Minute {
			return nil
		}
		_, err := c.Do("PING")
		return err
	},
}

func main() {
	fmt.Println("hello, i am sender.")

	c := pool.Get()
	defer c.Close()

	c.Do("PUBLISH", "my-first-channel", "hello, from sender.")
}
