package main

import (
	"bufio"
	"chat/sandbox/redis/04/db"
	"chat/sandbox/redis/04/element"
	"chat/sandbox/redis/04/receiver"
	"chat/sandbox/redis/04/redis"
	"chat/sandbox/redis/04/room"
	"chat/sandbox/redis/04/sender"
	"context"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: not enough args")
		os.Exit(1)
	}
	roomID := os.Args[1]

	rAddr := ":6379"
	rp := redis.NewPool(rAddr, 10*time.Minute)
	rr := receiver.NewRedisReceiver(rAddr)
	rs := sender.NewRedisSender(rp)
	rdb := db.NewDatabase(nil)
	sv := room.NewSupervisor(rr, rs, rdb)

	id, _ := strconv.ParseInt(roomID, 10, 64)
	room := sv.GetRoom(id)

	ctx := context.Background()
	go room.Run(ctx)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			return
		}

		text := scanner.Text()
		fields := strings.Fields(text)
		if len(fields) > 1 {
			user, _ := user.Current()
			msg := element.Message{User: user.Username, Body: fields[0]}
			room.Send(ctx, msg, fields[1:]...)
		} else {
			fmt.Println("ERROR: not enough args")
		}
	}
}

// --------------------
// Redis
// --------------------
