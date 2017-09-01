package models

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
)

type Session struct {
	ID        string
	CompanyID int
	ExpiresAt time.Time
}

const sessionSeed = "eerugb23j;o3ijbago4ih"
const sessionPeriod = 300

func registerSession(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Session{}, "sessions").SetKeys(false, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("CompanyID").Rename("company_id")
	t.ColMap("ExpiresAt").Rename("expires_at")
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
		return nil, errs.ErrNotFound
	}
	session := rows[0].(*Session)
	return session, nil
}

func CreateSession(company *Company) (*Session, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	session, err := CreateSessionInTx(company, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return session, nil
}

func CreateSessionInTx(company *Company, tx *gorp.Transaction) (*Session, error) {
	now := time.Now()

	seedStr := company.EmailAddress + sessionSeed + now.String()
	sessionID := fmt.Sprintf(
		"%x",
		sha256.Sum256([]byte(seedStr)),
	)

	session := &Session{
		ID:        sessionID,
		CompanyID: company.ID,
		ExpiresAt: now.Add(sessionPeriod),
	}
	err := tx.Insert(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
