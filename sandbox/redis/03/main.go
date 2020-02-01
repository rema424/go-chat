package main

import (
	"bufio"
	"chat/sandbox/terminal"
	"fmt"
	"os"
	"strings"
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
	go receiver()
	sender()
}

func receiver() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: not enough args")
		os.Exit(1)
	}

	channels := make([]interface{}, len(os.Args[1:]))
	for i, v := range os.Args[1:] {
		channels[i] = v
	}

	c := pool.Get()
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	defer psc.Close()

	if err := psc.Subscribe(channels...); err != nil {
		panic(err)
	}

	for {
		fmt.Println("waiting a notification")
		notification := psc.Receive()
		fmt.Println("catched a notification")
		switch n := notification.(type) {
		case error:
			fmt.Printf(terminal.Red("err: %q\n"), n)
		case redis.Pong:
			fmt.Printf(terminal.Magenta("pong: data: %s\n"), n.Data)
		case redis.Message:
			fmt.Printf(terminal.Green("msg: channel: %s, data: %s\n"), n.Channel, string(n.Data))
		case redis.Subscription:
			fmt.Printf(terminal.Blue("sub: kind: %s, channel: %s, count: %d\n"), n.Kind, n.Channel, n.Count)
		}
	}
}

func sender() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			return
		}

		text := scanner.Text()
		fields := strings.Fields(text)
		if len(fields) > 1 {
			msg := fields[0]
			channels := fields[1:]
			c := pool.Get()
			for _, ch := range channels {
				c.Do("PUBLISH", ch, msg)
			}
			c.Close()
		}
	}
}
