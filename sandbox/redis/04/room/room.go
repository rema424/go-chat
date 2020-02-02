package room

import (
	"chat/sandbox/redis/04/db"
	"chat/sandbox/redis/04/element"
	"chat/sandbox/redis/04/receiver"
	"chat/sandbox/redis/04/sender"
	"chat/sandbox/terminal"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

// --------------------
// SuperVisor
// --------------------

type Supervisor struct {
	cache   map[int64]*Room
	rceiver receiver.Receiver
	sender  sender.Sender
	databae db.Database
}

func NewSupervisor(
	r receiver.Receiver,
	s sender.Sender,
	d db.Database,
) *Supervisor {
	return &Supervisor{
		cache:   make(map[int64]*Room),
		rceiver: r,
		sender:  s,
		databae: d,
	}
}

func (sv *Supervisor) GetRoom(id int64) *Room {
	if v, ok := sv.cache[id]; ok {
		return v
	}
	room := NewRoom(id, sv.rceiver, sv.sender, sv.databae)
	sv.cache[id] = room
	return room
}

// --------------------
// Room
// --------------------

type Room struct {
	id      int64
	once    sync.Once
	rceiver receiver.Receiver
	sender  sender.Sender
	databae db.Database
}

func NewRoom(
	id int64,
	r receiver.Receiver,
	s sender.Sender,
	d db.Database,
) *Room {
	return &Room{
		id:      id,
		rceiver: r,
		sender:  s,
		databae: d,
	}
}

func (r *Room) Run(ctx context.Context) {
	r.once.Do(func() {
		var (
			onStart receiver.OnStart = func() error {
				fmt.Println("onStart called")
				return nil
			}
			onMessage receiver.OnMessage = func(channel string, data []byte) error {
				var msg element.Message
				_ = json.Unmarshal(data, &msg)
				fmt.Printf(terminal.Green("msg: channel: %s, data: %#v\n"), channel, msg)
				return nil
			}
		)

		if err := r.rceiver.Receive(ctx, onStart, onMessage, strconv.FormatInt(r.id, 10)); err != nil {
			panic(err)
		}
	})
}

func (r *Room) Send(ctx context.Context, msg element.Message, channels ...string) (element.Message, error) {
	msg, err := r.databae.Create(ctx, msg)
	if err != nil {
		return msg, err
	}
	b, _ := json.Marshal(msg)
	r.sender.Send(b, channels...)
	return msg, err
}
