package main

import (
	"bufio"
	"chat/sandbox/terminal"
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
	fmt.Println("hello, i am receiver.")
	c := pool.Get()
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	defer psc.Close()

	if err := psc.Subscribe("my-first-channel"); err != nil {
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

	// sendLoop(os.Stdin, os.Stdout)
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
