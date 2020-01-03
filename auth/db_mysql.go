package auth

import (
	"context"

	"github.com/rema424/sqlxx"
	"golang.org/x/xerrors"
)

func NewMySQL(db *sqlxx.Accessor) Database {
	return &MySQL{db}
}

type MySQL struct {
	db *sqlxx.Accessor
}

func (m *MySQL) RunInTx(
	ctx context.Context,
	txFn func(ctx context.Context) (interface{}, error),
) (v interface{}, err error, rberr error) {
	return m.db.RunInTx(ctx, txFn)
}

func (m *MySQL) CreateUser(ctx context.Context, s *Session) (*Session, error) {
	q := `INSERT INTO user (email, password) VALUES (:email, :password);`
	res, err := m.db.NamedExec(ctx, q, s.User)
	if err != nil {
		return s, xerrors.Errorf(": %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return s, xerrors.Errorf(": %w", err)
	}
	s.User.ID = id
	return s, nil
}

func (m *MySQL) CreateSession(ctx context.Context, s *Session) (*Session, error) {
	q := `INSERT INTO session (id, csrf, user_id, expire_at) VALUES (:id, :csrf, :user.id, :expire_at);`
	_, err := m.db.NamedExec(ctx, q, s)
	if err != nil {
		return s, xerrors.Errorf(": %w", err)
	}
	return s, nil
}

func (m *MySQL) GetUserByID(ctx context.Context, id int64) (*User, error) {
	q := `SELECT id, email, password FROM user WHERE id = ?;`
	var u User
	if err := m.db.Get(ctx, &u, q, id); err != nil {
		return &u, xerrors.Errorf(": %w", err)
	}
	return &u, nil
}

func (m *MySQL) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	q := `SELECT id, email, password FROM user WHERE email = ?;`
	var u User
	if err := m.db.Get(ctx, &u, q, email); err != nil {
		return &u, xerrors.Errorf(": %w", err)
	}
	return &u, nil
}

func (m *MySQL) GetSessionByID(ctx context.Context, id string) (*Session, error) {
	q := `
  SELECT
    s.id,
    s.csrf,
    s.expire_at,
    u.id AS 'user.id',
    u.email AS 'user.email',
    u.password AS 'user.password'
  FROM session AS s
  INNER JOIN user AS u ON u.id = s.user_id
  WHERE s.id = ?;
  `
	var s Session
	if err := m.db.Get(ctx, &s, q, id); err != nil {
		return &s, xerrors.Errorf(": %w", err)
	}
	return &s, nil
}
