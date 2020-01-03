package auth

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/xerrors"
)

var (
	ErrGenerateHash = xerrors.New("auth: password: failed to generate hash")
)

func NewSession(u *User) *Session {
	if u == nil {
		u = &User{}
	}
	return &Session{
		ID:       uuid.New().String(),
		CSRF:     uuid.New().String(),
		ExpireAt: time.Now().AddDate(0, 2, 0).Unix(),
		User:     u,
	}
}

func NewUser(email, row string) (*User, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(row), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, xerrors.Errorf("%s: %w", ErrGenerateHash, err)
	}
	return &User{Email: email, Password: string(b)}, nil
}

type Session struct {
	ID       string `db:"id"`
	CSRF     string `db:"csrf"`
	ExpireAt int64  `db:"expire_at"`
	User     *User  `db:"user"`
}

type User struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (u *User) VerifyPassword(row string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(row))
}
