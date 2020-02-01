package auth

import "context"

type FakeDB struct {
	Database
	FakeRunInTx        func(context.Context, func(context.Context) (interface{}, error)) (v interface{}, err error, rberr error)
	FakeCreateUser     func(context.Context, *Session) (*Session, error)
	FakeCreateSession  func(context.Context, *Session) (*Session, error)
	FakeGetUserByID    func(context.Context, int64) (*User, error)
	FakeGetUserByEmail func(context.Context, string) (*User, error)
	FakeGetSessionByID func(context.Context, string) (*Session, error)
}

func (fd *FakeDB) RunInTx(ctx context.Context, txFn func(context.Context) (interface{}, error)) (v interface{}, err error, rberr error) {
	return fd.FakeRunInTx(ctx, txFn)
}

// func fakeRunInTx(ctx context.Context, txFn func(context.Context) (interface{}, error)) (v interface{}, err error, rberr error) {
// 	return
// }

func (fd *FakeDB) CreateUser(ctx context.Context, s *Session) (*Session, error) {
	return fd.FakeCreateUser(ctx, s)
}

func (fd *FakeDB) CreateSession(ctx context.Context, s *Session) (*Session, error) {
	return fd.FakeCreateSession(ctx, s)
}

func (fd *FakeDB) GetUserByID(ctx context.Context, id int64) (*User, error) {
	return fd.FakeGetUserByID(ctx, id)
}

func (fd *FakeDB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return fd.FakeGetUserByEmail(ctx, email)
}

func (fd *FakeDB) GetSessionByID(ctx context.Context, id string) (*Session, error) {
	return fd.FakeGetSessionByID(ctx, id)
}
