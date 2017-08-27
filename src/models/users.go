package models

import (
	"crypto/sha256"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
)

type User struct {
	ID           int
	Name         string
	EmailAddress string
	Password     string
}

const passwordSeed = "g3irugb34trsajdkgh"

func registerUser(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(User{}, "users").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("Name").Rename("name")
	t.ColMap("EmailAddress").Rename("emailAddress")
	t.ColMap("Password").Rename("password")
}

func comparePassword(rawPassword string) bool {
	return rawPassword == generateUserPasswordHash(rawPassword)
}

func generateUserPasswordHash(rawPassword string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(rawPassword+passwordSeed)))
}

func FindUser(userID int) (*User, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(User{}, "select * from users where id = ?", userID)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}
	user := rows[0].(*User)
	return user, nil
}

func FindUserByEmailAddress(emailAddress string) (*User, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(User{}, "select * from users where emailAddress = ?", emailAddress)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}
	user := rows[0].(*User)
	return user, nil
}

func CreateUser(name string, emailAddress string, password string) (*User, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	user, err := CreateUserInTx(name, emailAddress, password, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return user, nil
}

func CreateUserInTx(name string, emailAddress string, rawPassword string, tx *gorp.Transaction) (*User, error) {
	user := &User{
		Name:         name,
		EmailAddress: emailAddress,
		Password:     generateUserPasswordHash(rawPassword),
	}
	err := tx.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
