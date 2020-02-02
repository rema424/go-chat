package sender

import "github.com/gomodule/redigo/redis"

// --------------------
// Sender
// --------------------

type Sender interface {
	Send(data []byte, channels ...string)
}

// --------------------
// RedisSender
// --------------------

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
