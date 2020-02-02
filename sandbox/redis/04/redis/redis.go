package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

const healthCheckPeriod = time.Minute

func NewPool(addr string, healthCheckPeriod time.Duration) *redis.Pool {
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
