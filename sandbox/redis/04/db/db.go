package db

import (
	"context"
	"fmt"

	"chat/sandbox/redis/04/element"

	"github.com/rema424/sqlxx"
)

type Database interface {
	Create(ctx context.Context, msg element.Message) (element.Message, error)
}

type database struct {
	db *sqlxx.DB
}

func NewDatabase(db *sqlxx.DB) Database {
	return &database{db: db}
}

func (db *database) Create(ctx context.Context, msg element.Message) (element.Message, error) {
	if db.db == nil {
		fmt.Println("plz implement db")
		return msg, nil
	}
	q := `INSERT INTO message (user, body) VALUES (:user, :body);`
	res, err := db.db.NamedExec(ctx, q, msg)
	if err != nil {
		return msg, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return msg, err
	}
	msg.ID = id
	return msg, nil
}
