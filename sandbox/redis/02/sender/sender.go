package main

import (
	"fmt"
	"os"
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
	if len(os.Args) < 3 {
		fmt.Println("ERROR: not enough args")
		os.Exit(1)
	}

	msg := os.Args[1]
	channels := os.Args[2:]

	fmt.Println("hello, i am sender.")

	c := pool.Get()
	defer c.Close()

	for _, ch := range channels {
		c.Do("PUBLISH", ch, msg)
	}
}
