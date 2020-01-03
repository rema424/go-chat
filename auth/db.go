package auth

import "context"

type Database interface {
	RunInTx(context.Context, func(context.Context) (interface{}, error)) (interface{}, error, error)
	CreateUser(ctx context.Context, s *Session) (*Session, error)
	CreateSession(ctx context.Context, s *Session) (*Session, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
