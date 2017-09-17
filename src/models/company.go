package models

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
)

type Company struct {
	ID           int
	Name         string
	EmailAddress string
	Password     string
	CreatedAt    time.Time
	ImageLink    string
}

func registerCompany(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Company{}, "companies").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("Name").Rename("name")
	t.ColMap("EmailAddress").Rename("email_address")
	t.ColMap("Password").Rename("password")
	t.ColMap("CreatedAt").Rename("created_at")
	t.ColMap("ImageLink").Rename("image_link")
}

func (self *Company) ComparePassword(rawPassword string) bool {
	return generatePasswordHash(rawPassword) == self.Password
}

func generatePasswordHash(rawPassword string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(rawPassword+"wogiba;eworugba;w")))
}

func FindCompany(companyID int) (*Company, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Company{}, "select * from companies where id = ?", companyID)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}
	company := rows[0].(*Company)
	return company, nil
}

func FindCompanyByEmailAddress(emailAddress string) (*Company, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Company{}, "select * from companies where email_address = ?", emailAddress)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}
	company := rows[0].(*Company)
	return company, nil
}

func CreateCompany(name string, emailAddress string, password string) (*Company, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	company, err := CreateCompanyInTx(name, emailAddress, password, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return company, nil
}

func CreateCompanyInTx(name string, emailAddress string, rawPassword string, tx *gorp.Transaction) (*Company, error) {
	company := &Company{
		Name:         name,
		EmailAddress: emailAddress,
		Password:     generatePasswordHash(rawPassword),
		CreatedAt:    time.Now(),
		ImageLink:    "",
	}
	err := tx.Insert(company)
	if err != nil {
		return nil, err
	}

	return company, nil
}

func (self *Company) SetPassword(rawPassword string) {
	self.Password = generatePasswordHash(rawPassword)
}

func (self *Company) Update() error {
	db, err := Db()
	if err != nil {
		return err
	}

	tx, _ := db.Begin()
	_, err = tx.Update(self)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
