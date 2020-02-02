package main

import (
	"bufio"
	"chat/sandbox/terminal"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

const healthCheckPeriod = time.Minute

var pool = NewRedisPool(":6379", healthCheckPeriod)

type Message struct {
	User string `json:"user"`
	Body string `json:"body"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: not enough args")
		os.Exit(1)
	}
	ctx := context.Background()
	rr := NewRedisReceiver(":6379")
	var (
		onStart OnStart = func() error {
			fmt.Println("onStart called")
			return nil
		}
		onMessage OnMessage = func(channel string, data []byte) error {
			var msg Message
			_ = json.Unmarshal(data, &msg)
			fmt.Printf(terminal.Green("msg: channel: %s, data: %#v\n"), channel, msg)
			return nil
		}
	)
	go func() {
		if err := rr.Receive(ctx, onStart, onMessage, os.Args[1:]...); err != nil {
			panic(err)
		}
	}()

	rs := NewRedisSender(pool)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			return
		}

		text := scanner.Text()
		fields := strings.Fields(text)
		if len(fields) > 1 {
			user, _ := user.Current()
			msg := Message{User: user.Username, Body: fields[0]}
			b, _ := json.Marshal(msg)
			rs.Send(b, fields[1:]...)
		} else {
			fmt.Println("ERROR: not enough args")
		}
	}
}

// --------------------
// Redis
// --------------------

func NewRedisPool(addr string, healthCheckPeriod time.Duration) *redis.Pool {
	p := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				addr,
				redis.DialReadTimeout(10*time.Second+healthCheckPeriod),
				redis.DialWriteTimeout(10*time.Second),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < healthCheckPeriod {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         5,
		MaxActive:       0,
		IdleTimeout:     5 * time.Minute,
		Wait:            false,
		MaxConnLifetime: 0,
	}

	c := p.Get()
	if _, err := c.Do("PING"); err != nil {
		panic(err)
	}
	c.Close()
	return p
}

// --------------------
// Receiver
// --------------------

type (
	OnStart   func() error
	OnMessage func(channel string, data []byte) error
)

type Receiver interface {
	Receive(ctx context.Context, onStart OnStart, onMessage OnMessage, channels ...string) error
}

type RedisReceiver struct {
	addr string
}

func NewRedisReceiver(addr string) Receiver {
	return &RedisReceiver{addr: addr}
}

func (rr *RedisReceiver) Receive(ctx context.Context, onStart OnStart, onMessage OnMessage, channels ...string) error {
	c, err := redis.Dial(
		"tcp",
		rr.addr,
		redis.DialReadTimeout(10*time.Second+healthCheckPeriod),
		redis.DialWriteTimeout(10*time.Second),
	)
	if err != nil {
		return err
	}
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}

	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return err
	}

	done := make(chan error, 1)

	go func() {
		for {
			fmt.Println("waiting")
			notification := psc.Receive()
			fmt.Println("catched a notification")

			switch n := notification.(type) {
			case error:
				fmt.Printf(terminal.Red("err: %q\n"), n)
				done <- n // goroutine では戻り値の代わりに channel で値を渡す
				return
			case redis.Pong:
				fmt.Printf(terminal.Magenta("pong: data: %s\n"), n.Data)
			case redis.Message:
				if err := onMessage(n.Channel, n.Data); err != nil {
					done <- err
					return // goroutine では戻り値の代わりに channel で値を渡す
				}
			case redis.Subscription: // subscribe or unsubscribe
				fmt.Printf(terminal.Blue("sub: kind: %s, channel: %s, count: %d\n"), n.Kind, n.Channel, n.Count)
				switch n.Count {
				case len(channels): // Subscribe
					if err := onStart(); err != nil {
						done <- err
						return
					}
				case 0: // Unsubscribe
					done <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()

loop:
	for err == nil {
		select {
		case err := <-done:
			return err
		case <-ticker.C:
			err = psc.Ping("")
		case <-ctx.Done():
			break loop
		}
	}

	psc.Unsubscribe()
	return <-done
}

// --------------------
// Sender
// --------------------

type Sender interface {
	Send(data []byte, channels ...string)
}

type RedisSender struct {
	pool *redis.Pool
}

func NewRedisSender(p *redis.Pool) Sender {
	return &RedisSender{pool: p}
}

func (rs *RedisSender) Send(data []byte, channels ...string) {
	c := rs.pool.Get()
	for _, ch := range channels {
		c.Do("PUBLISH", ch, data)
	}
	c.Close()
}
