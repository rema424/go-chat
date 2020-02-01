package broker

// import (
//   "github.com/gomodule/redigo/redis"
// )

// type Broker interface {
// 	// Close() error
// 	// Publish(v interface{}) error
// 	// Recieve() (interface{}, error)
// }

// var pool = &redis.Pool{
//   MaxIdle:     3,
//   IdleTimeout: 240 * time.Second,
//   Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", ":6379") },
//   TestOnBorrow: func(c redis.Conn, t time.Time) error {
//     if time.Since(t) < time.Minute {
//       return nil
//     }
//     _, err := c.Do("PING")
//     return err
//   },
// }

// type RedisPubSub struct{
//   c *redis.Conn
// }

// func NewRedis() Broker {
// 	return &Redis{}
// }

// // func (r *Redis) Close() error {
// //   return
// // }
