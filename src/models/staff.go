package models

import (
	"time"

	"github.com/ren-motomura/lesson-manager-api-server/src/errs"

	"github.com/go-gorp/gorp"
)

type Staff struct {
	ID        int
	Name      string
	CompanyID int
	CreatedAt time.Time
	IsValid   bool
	ImageLink string
}

func registerStaff(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Staff{}, "staffs").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("Name").Rename("name")
	t.ColMap("CompanyID").Rename("company_id")
	t.ColMap("CreatedAt").Rename("created_at")
	t.ColMap("IsValid").Rename("is_valid")
	t.ColMap("ImageLink").Rename("image_link")
}

func (self *Staff) Delete() error {
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

func (self *Staff) DeleteInTx(tx *gorp.Transaction) error {
	self.IsValid = false
	_, err := tx.Update(self)
	if err != nil {
		return err
	}

	return nil
}

func FindStaff(id int, forUpdate bool, tx *gorp.Transaction) (*Staff, error) {
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

	rows, err := selector.Select(Staff{}, "select * from staffs where id = ? and is_valid = ? "+forUpStatement, id, true)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}
	staff := rows[0].(*Staff)
	return staff, nil
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

func CreateStaff(name string, imageLink string, company *Company) (*Staff, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	staff, err := CreateStaffInTx(name, imageLink, company, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return staff, nil
}

func CreateStaffInTx(name string, imageLink string, company *Company, tx *gorp.Transaction) (*Staff, error) {
	staff := &Staff{
		Name:      name,
		CompanyID: company.ID,
		CreatedAt: time.Now(),
		IsValid:   true,
		ImageLink: imageLink,
	}
	err := tx.Insert(staff)
	if err != nil {
		return nil, err
	}

	return staff, nil
}
