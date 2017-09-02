package models

import (
	"time"

	"github.com/go-gorp/gorp"
)

type Studio struct {
	ID        int
	Name      string
	CompanyID int
	CreatedAt time.Time
	IsValid   bool
}

func registerStudio(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Studio{}, "studios").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("Name").Rename("name")
	t.ColMap("CompanyID").Rename("company_id")
	t.ColMap("CreatedAt").Rename("created_at")
	t.ColMap("IsValid").Rename("is_valid")
}

func (self *Studio) Delete() error {
	db, err := Db()
	if err != nil {
		return err
	}

	self.IsValid = false
	_, err = db.Update(self)
	if err != nil {
		return err
	}

	return nil
}

func SelectStudiosByCompany(company *Company) ([]*Studio, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Studio{}, "select * from studios where company_id = ? and is_valid = ?", company.ID, true)
	if err != nil {
		return nil, err
	}
	studios := make([]*Studio, len(rows))
	for i, row := range rows {
		studios[i] = row.(*Studio)
	}
	return studios, nil
}

func CreateStudio(name string, company *Company) (*Studio, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	studio, err := CreateStudioInTx(name, company, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return studio, nil
}

func CreateStudioInTx(name string, company *Company, tx *gorp.Transaction) (*Studio, error) {
	studio := &Studio{
		Name:      name,
		CompanyID: company.ID,
		CreatedAt: time.Now(),
		IsValid:   true,
	}
	err := tx.Insert(studio)
	if err != nil {
		return nil, err
	}

	return studio, nil
}
