package auth

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSignup(t *testing.T) {
	tests := []struct {
		in      *User
		out     *Session
		wantErr bool
		fakeDB  Database
	}{
		{
			in: &User{},
			fakeDB: &FakeDB{
				FakeRunInTx: func(ctx context.Context, txFn func(context.Context) (interface{}, error)) (interface{}, error, error) {
					v, err := txFn(ctx)
					return v, err, nil
				},
				FakeCreateUser: func(ctx context.Context, s *Session) (*Session, error) {
					s.User.ID = 10
					return s, nil
				},
				FakeCreateSession: func(ctx context.Context, s *Session) (*Session, error) { return s, nil },
			},
			out: &Session{User: &User{ID: 10}},
		},
	}

	ii := &InteractorImpl{}
	for _, tt := range tests {
		ii.db = tt.fakeDB
		ctx := context.Background()
		s, err := ii.Signup(ctx, tt.in)
		if tt.wantErr && err == nil {
			t.Error("want non-nil error")
		}
		t.Logf("%#v, %#v", s.User, tt.in)
		if diff := cmp.Diff(s.User, tt.in); diff != "" {
			t.Errorf("wrong user\n%s", diff)
		}
	}
}
