package models

import (
	"time"

	"github.com/go-gorp/gorp"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
)

type Customer struct {
	ID          int
	Name        string
	Description string
	CompanyID   int
	CreatedAt   time.Time
	IsValid     bool
}

func registerCustomer(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Customer{}, "customers").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("Name").Rename("name")
	t.ColMap("Description").Rename("description")
	t.ColMap("CompanyID").Rename("company_id")
	t.ColMap("CreatedAt").Rename("created_at")
	t.ColMap("IsValid").Rename("is_valid")
}

func (self *Customer) Delete() error {
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

func (self *Customer) DeleteInTx(tx *gorp.Transaction) error {
	self.IsValid = false
	_, err := tx.Update(self)
	if err != nil {
		return err
	}

	return nil
}

func FindCustomer(id int, forUpdate bool, tx *gorp.Transaction) (*Customer, error) {
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

	rows, err := selector.Select(Customer{}, "select * from customers where id = ? and is_valid = ? "+forUpStatement, id, true)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}
	customer := rows[0].(*Customer)
	return customer, nil
}

func SelectCustomersByCompany(company *Company) ([]*Customer, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Customer{}, "select * from customers where company_id = ? and is_valid = ?", company.ID, true)
	if err != nil {
		return nil, err
	}
	customers := make([]*Customer, len(rows))
	for i, row := range rows {
		customers[i] = row.(*Customer)
	}
	return customers, nil
}

func CreateCustomer(name string, description string, company *Company) (*Customer, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	customer, err := CreateCustomerInTx(name, description, company, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return customer, nil
}

func CreateCustomerInTx(name string, description string, company *Company, tx *gorp.Transaction) (*Customer, error) {
	customer := &Customer{
		Name:        name,
		Description: description,
		CompanyID:   company.ID,
		CreatedAt:   time.Now(),
		IsValid:     true,
	}
	err := tx.Insert(customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (self *Customer) UpdateInTx(tx *gorp.Transaction) error {
	_, err := tx.Update(self)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
