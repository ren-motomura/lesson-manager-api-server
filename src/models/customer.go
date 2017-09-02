package models

import (
	"errors"
	"time"

	"github.com/ren-motomura/lesson-manager-api-server/src/errs"

	"github.com/go-gorp/gorp"
)

type Customer struct {
	ID          int
	Name        string
	Description string
	CompanyID   int
	CardID      string
	Credit      int
	CreatedAt   time.Time
	IsValid     bool
}

func registerCustomer(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Customer{}, "customers").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("Name").Rename("name")
	t.ColMap("Description").Rename("description")
	t.ColMap("CompanyID").Rename("company_id")
	t.ColMap("CardID").Rename("card_id")
	t.ColMap("Credit").Rename("credit")
	t.ColMap("CreatedAt").Rename("created_at")
	t.ColMap("IsValid").Rename("is_valid")
}

func (self *Customer) Delete() error {
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

func FindCustomerByCompanyAndCardID(company *Company, cardID string) (*Customer, error) {
	if cardID == "" {
		return nil, errors.New("invalid card ID")
	}

	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Customer{}, "select * from customers where company_id = ? and card_id = ? and is_valid = ?", company.ID, cardID, true)
	if err != nil {
		return nil, err
	}

	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}

	customer := rows[0].(*Customer)
	return customer, nil
}

func CreateCustomer(name string, description string, company *Company, cardID string, credit int) (*Customer, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	customer, err := CreateCustomerInTx(name, description, company, cardID, credit, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return customer, nil
}

func CreateCustomerInTx(name string, description string, company *Company, cardID string, credit int, tx *gorp.Transaction) (*Customer, error) {
	customer := &Customer{
		Name:        name,
		Description: description,
		CompanyID:   company.ID,
		CardID:      cardID,
		Credit:      credit,
		CreatedAt:   time.Now(),
		IsValid:     true,
	}
	err := tx.Insert(customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
