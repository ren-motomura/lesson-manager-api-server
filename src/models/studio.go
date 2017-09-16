package models

import (
	"time"

	"github.com/go-gorp/gorp"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
)

type Studio struct {
	ID          int
	Name        string
	Address     string
	PhoneNumber string
	CompanyID   int
	CreatedAt   time.Time
	IsValid     bool
	ImageLink   string
}

func registerStudio(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Studio{}, "studios").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("Name").Rename("name")
	t.ColMap("Address").Rename("address")
	t.ColMap("PhoneNumber").Rename("phone_number")
	t.ColMap("CompanyID").Rename("company_id")
	t.ColMap("CreatedAt").Rename("created_at")
	t.ColMap("IsValid").Rename("is_valid")
	t.ColMap("ImageLink").Rename("image_link")
}

func (self *Studio) Delete() error {
	db, err := Db()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = self.DeleteInTx(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (self *Studio) DeleteInTx(tx *gorp.Transaction) error {
	self.IsValid = false
	_, err := tx.Update(self)
	if err != nil {
		return err
	}

	return nil
}

func FindStudio(id int, forUpdate bool, tx *gorp.Transaction) (*Studio, error) {
	forUpStatement := ""
	if forUpdate {
		forUpStatement = "for update"
	}

	var selector Selector
	if tx != nil {
		selector = tx
	} else {
		db, err := Db()
		if err != nil {
			return nil, err
		}
		selector = db
	}

	rows, err := selector.Select(Studio{}, "select * from studios where id = ? and is_valid = ? "+forUpStatement, id, true)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}
	studio := rows[0].(*Studio)
	return studio, nil
}

func FindStudioByCompanyAndName(company *Company, name string) (*Studio, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Studio{}, "select * from studios where company_id = ? and name = ?", company.ID, name)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}
	studio := rows[0].(*Studio)
	return studio, nil
}

func SelectStudiosByCompany(company *Company) ([]*Studio, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Studio{}, "select * from studios where company_id = ? and is_valid = ? order by created_at asc", company.ID, true)
	if err != nil {
		return nil, err
	}
	studios := make([]*Studio, len(rows))
	for i, row := range rows {
		studios[i] = row.(*Studio)
	}
	return studios, nil
}

func CreateStudio(name string, address string, phoneNumber string, company *Company, imageLink string) (*Studio, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	studio, err := CreateStudioInTx(name, address, phoneNumber, company, imageLink, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return studio, nil
}

func CreateStudioInTx(name string, address string, phoneNumber string, company *Company, imageLink string, tx *gorp.Transaction) (*Studio, error) {
	studio := &Studio{
		Name:        name,
		Address:     address,
		PhoneNumber: phoneNumber,
		CompanyID:   company.ID,
		CreatedAt:   time.Now(),
		IsValid:     true,
		ImageLink:   imageLink,
	}
	err := tx.Insert(studio)
	if err != nil {
		return nil, err
	}

	return studio, nil
}
