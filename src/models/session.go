package models

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/go-gorp/gorp"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

const sessionSeed = "eerugb23j;o3ijbago4ih"
const sessionPeriod = 300

func registerSession(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Session{}, "sessions").SetKeys(false, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("UserID").Rename("userId")
	t.ColMap("ExpiresAt").Rename("expiresAt")
}

func FindSession(sessionID string) (*Session, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Session{}, "select * from sessions where id = ?", sessionID)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errors.New("session not found")
	}
	session := rows[0].(*Session)
	return session, nil
}

func CreateSession(user *User) (*Session, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	session, err := CreateSessionInTx(user, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return session, nil
}

func CreateSessionInTx(user *User, tx *gorp.Transaction) (*Session, error) {
	now := time.Now()

	seedStr := user.EmailAddress + sessionSeed + now.String()
	sessionID := fmt.Sprintf(
		"%x",
		sha256.Sum256([]byte(seedStr)),
	)

	session := &Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: now.Add(sessionPeriod),
	}
	err := tx.Insert(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
