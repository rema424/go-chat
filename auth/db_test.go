/*
  mysql.server start
  mysql -uroot -e 'create database if not exists chattest;'
  mysql -uroot -e 'create user if not exists dbtester@localhost identified by "Passw0rd!";'
  mysql -uroot -e 'grant all privileges on chattest.* to dbtester@localhost;'
  mysql -uroot -e 'show databases;'
  mysql -uroot -e 'select host, user from mysql.user;'
  mysql -uroot -e 'show grants for dbtester@localhost;'
*/

package auth

import (
	"context"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/rema424/sqlxx"
)

const CreateUser = `
create table if not exists user (
  id bigint not null auto_increment,
  email varchar(255) character set latin1 collate latin1_bin not null default '',
  password varchar(255) not null default '',
  primary key (id),
  unique key (email)
);
`
const DropUser = `
drop table if exists user;
`

const CreateSession = `
create table if not exists session (
  id varchar(255) character set latin1 collate latin1_bin not null default '',
  csrf varchar(255) not null default '',
  user_id bigint not null default 0,
  expire_at bigint not null default 0,
  primary key (id),
  foreign key (user_id) references user (id) on delete cascade on update cascade,
  key (user_id)
);
`
const DropSession = `
drop table if exists session;
`

func TestMySQL(t *testing.T) {
	dbx, err := sqlx.Connect("mysql", "dbtester:Passw0rd!@tcp(127.0.0.1:3306)/chattest?collation=utf8mb4_bin&interpolateParams=true&parseTime=true&maxAllowedPacket=0")
	if err != nil {
		log.Fatalf("sqlx.Connect: %+v", err)
	}
	dbx.MustExec(DropSession)
	dbx.MustExec(DropUser)
	dbx.MustExec(CreateUser)
	dbx.MustExec(CreateSession)

	db, err := sqlxx.Open(dbx)
	if err != nil {
		log.Fatalf("sqlx.Connect: %+v", err)
	}

	testDB(t, NewMySQL(db))
}

func testDB(t *testing.T, db Database) {
	t.Helper()

	// --------------------
	// setup
	// --------------------
	u, err := NewUser("abc@example.com", "Passw0rd!")
	if err != nil {
		t.Fatalf("NewUser returned error: %+v", err)
	}
	s := NewSession(u)
	ctx := context.Background()

	// --------------------
	// CreateUser
	// --------------------
	s, err = db.CreateUser(ctx, s)
	if err != nil {
		t.Fatalf("CreateUser returned error: %+v", err)
	}
	if s.User.ID == 0 {
		t.Fatalf("CreateUser did not return user id: got %d", s.User.ID)
	}
	_, err = db.CreateUser(ctx, s)
	if err == nil {
		t.Fatalf("CreateUser want non-nil error")
	} else {
		t.Logf("%+v", err)
	}

	// --------------------
	// CreateSession
	// --------------------
	s, err = db.CreateSession(ctx, s)
	if err != nil {
		t.Fatalf("CreateSession returned error: %+v", err)
	}
	_, err = db.CreateSession(ctx, s)
	if err == nil {
		t.Fatalf("CreateSession want non-nil error")
	} else {
		t.Logf("%+v", err)
	}

	// --------------------
	// GetUserByID
	// --------------------
	gotUserByID, err := db.GetUserByID(ctx, s.User.ID)
	if err != nil {
		t.Fatalf("GetUserByID returned error: %+v", err)
	}
	if diff := cmp.Diff(gotUserByID, s.User); diff != "" {
		t.Fatalf("GetUserByID returned wrong result: \n%s", diff)
	}
	_, err = db.GetUserByID(ctx, 9999)
	if err == nil {
		t.Fatalf("GetUserByID want non-nil error")
	} else {
		t.Logf("%+v", err)
	}

	// --------------------
	// GetUserByEmail
	// --------------------
	gotUserByEmail, err := db.GetUserByEmail(ctx, s.User.Email)
	if err != nil {
		t.Fatalf("GetUserByEmail returned error: %+v", err)
	}
	if diff := cmp.Diff(gotUserByEmail, s.User); diff != "" {
		t.Fatalf("GetUserByEmail returned wrong result: \n%s", diff)
	}
	_, err = db.GetUserByEmail(ctx, "abcdef")
	if err == nil {
		t.Fatalf("GetUserByEmail want non-nil error")
	} else {
		t.Logf("%+v", err)
	}

	// --------------------
	// GetSessionByID
	// --------------------
	gotSessionByID, err := db.GetSessionByID(ctx, s.ID)
	if err != nil {
		t.Fatalf("GetSessionByID returned error: %+v", err)
	}
	if diff := cmp.Diff(gotSessionByID, s); diff != "" {
		t.Fatalf("GetSessionByID returned wrong result: \n%s", diff)
	}
	_, err = db.GetSessionByID(ctx, "abcdef")
	if err == nil {
		t.Fatalf("GetSessionByID want non-nil error")
	} else {
		t.Logf("%+v", err)
	}

	// --------------------
	// RunInTx
	// --------------------
	switch db := db.(type) {
	case *MySQL:
		testMySQLRunInTx(t, db)
	}
}

func testMySQLRunInTx(t *testing.T, db *MySQL) {
	t.Helper()

	// --------------------
	// setup
	// --------------------
	email1 := "mysql-tx-test-1@example.com"
	email2 := "mysql-tx-test-2@example.com"
	ctx := context.Background()

	// ----------------------------------------
	// transaction 1（success）
	// ----------------------------------------
	txFn1 := func(ctx context.Context) (interface{}, error) {
		u, err := NewUser(email1, "Passw0rd!")
		if err != nil {
			t.Fatal(err)
		}
		s := NewSession(u)

		s, err = db.CreateUser(ctx, s)
		if err != nil {
			return s, err
		}
		return db.CreateSession(ctx, s)
	}
	v, err, rberr := db.RunInTx(ctx, txFn1)
	if rberr != nil {
		t.Fatalf("RunInTx returned rollback error: %+v", rberr)
	}
	if err != nil {
		t.Fatalf("RunInTx returned error: %+v", err)
	}
	s, ok := v.(*Session)
	if !ok {
		t.Fatalf("RunInTx result is not *Session. got %T", v)
	}
	registeredSessionID := s.ID

	// ----------------------------------------
	// transaction 2（failuer and rollback）
	// ----------------------------------------
	txFn2 := func(ctx context.Context) (interface{}, error) {
		u, err := NewUser(email2, "Passw0rd!")
		if err != nil {
			t.Fatal(err)
		}
		s := NewSession(u)
		s.ID = registeredSessionID // trigger duplicate entry error

		s, err = db.CreateUser(ctx, s)
		if err != nil {
			return s, err
		}
		return db.CreateSession(ctx, s) // expect to occur duplicate entry error
	}
	_, err, rberr = db.RunInTx(ctx, txFn2)
	if rberr != nil {
		t.Fatalf("RunInTx returned rollback error: %+v", rberr)
	}
	if err == nil {
		t.Fatalf("RunInTx want non-nil error")
	} else {
		t.Logf("MySQL.RunInTx: transaction2:\n%+v", err)
	}

	// ----------------------------------------
	// rollback check
	// ----------------------------------------
	_, err = db.GetUserByEmail(ctx, email2)
	if err == nil {
		t.Fatalf("want non-nil error. RunInTx failed to rollback.")
	} else {
		t.Logf("MySQL.RunInTx: rollback check:\n%+v", err)
	}
}
