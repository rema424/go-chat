package auth

import (
	"context"
	"log"

	"golang.org/x/xerrors"
)

func NewInteractorImpl(db Database) Interactor {
	return &InteractorImpl{db}
}

type InteractorImpl struct {
	db Database
}

func (ii *InteractorImpl) Signup(ctx context.Context, u *User) (*Session, error) {
	s := NewSession(u)
	txFn := func(ctx context.Context) (interface{}, error) {
		var err error
		s, err = ii.db.CreateUser(ctx, s)
		if err != nil {
			return s, err
		}
		return ii.db.CreateSession(ctx, s)
	}
	v, err, rberr := ii.db.RunInTx(ctx, txFn)
	if rberr != nil {
		log.Printf("%+v", rberr)
	}
	if err != nil {
		return nil, err
	}
	s, ok := v.(*Session)
	if !ok {
		return nil, xerrors.New("auth: v not *Session")
	}
	return s, nil
}
