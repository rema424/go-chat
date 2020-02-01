package auth

import "context"

type Interactor interface {
	Signup(context.Context, *User) (*Session, error)
}
