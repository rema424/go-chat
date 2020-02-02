package receiver

import (
	"context"
	"fmt"
	"time"

	"chat/sandbox/terminal"

	"github.com/gomodule/redigo/redis"
)

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

// --------------------
// RedisReceiver
// --------------------

const healthCheckPeriod = time.Minute

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
			switch n := psc.Receive().(type) {
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
