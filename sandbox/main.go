package main

import (
	"bufio"
	"fmt"
	"io"
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
	c := pool.Get()
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	defer psc.Close()

	if err := psc.Subscribe("test"); err != nil {
		panic(err)
	}

	for {

	}

	sendLoop(os.Stdin, os.Stdout)
}

const prompt = ">> "

func sendLoop(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		if !scanner.Scan() { // blocking
			return
		}

		text := scanner.Text()

		io.WriteString(out, "\u001b[33m"+text+"\u001b[0m")
		io.WriteString(out, "\n")
	}
}
