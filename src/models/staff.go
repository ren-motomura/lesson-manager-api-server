package models

import (
	"time"

	"github.com/go-gorp/gorp"
)

type Staff struct {
	ID        int
	Name      string
	CompanyID int
	CreatedAt time.Time
	IsValid   bool
}

func registerStaff(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Staff{}, "staffs").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("Name").Rename("name")
	t.ColMap("CompanyID").Rename("company_id")
	t.ColMap("CreatedAt").Rename("created_at")
	t.ColMap("IsValid").Rename("is_valid")
}

func (self *Staff) Delete() error {
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

func SelectStaffsByCompany(company *Company) ([]*Staff, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Staff{}, "select * from staffs where company_id = ? and is_valid = ?", company.ID, true)
	if err != nil {
		return nil, err
	}
	staffs := make([]*Staff, len(rows))
	for i, row := range rows {
		staffs[i] = row.(*Staff)
	}
	return staffs, nil
}

func CreateStaff(name string, company *Company) (*Staff, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	staff, err := CreateStaffInTx(name, company, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return staff, nil
}

func CreateStaffInTx(name string, company *Company, tx *gorp.Transaction) (*Staff, error) {
	staff := &Staff{
		Name:      name,
		CompanyID: company.ID,
		CreatedAt: time.Now(),
		IsValid:   true,
	}
	err := tx.Insert(staff)
	if err != nil {
		return nil, err
	}

	return staff, nil
}
