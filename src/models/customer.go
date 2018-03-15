package models

import (
	"fmt"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
)

type Customer struct {
	ID           int
	CompanyID    int
	Name         string
	Kana         string
	IsValid      bool
	Birthday     int64
	Gender       Gender
	PostalCode1  string
	PostalCode2  string
	Address      string
	PhoneNumber  string
	JoinDate     int64
	EmailAddress string
	CanMail      bool
	CanEmail     bool
	CanCall      bool
	Description  string
	CreatedAt    int64
}

func registerCustomer(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Customer{}, "customers").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("CompanyID").Rename("company_id")
	t.ColMap("Name").Rename("name")
	t.ColMap("Kana").Rename("kana")
	t.ColMap("IsValid").Rename("is_valid")
	t.ColMap("Birthday").Rename("birthday")
	t.ColMap("Gender").Rename("gender")
	t.ColMap("PostalCode1").Rename("postal_code1")
	t.ColMap("PostalCode2").Rename("postal_code2")
	t.ColMap("Address").Rename("address")
	t.ColMap("PhoneNumber").Rename("phone_number")
	t.ColMap("JoinDate").Rename("join_date")
	t.ColMap("EmailAddress").Rename("email_address")
	t.ColMap("CanMail").Rename("can_mail")
	t.ColMap("CanEmail").Rename("can_email")
	t.ColMap("CanCall").Rename("can_call")
	t.ColMap("Description").Rename("description")
	t.ColMap("CreatedAt").Rename("created_at")
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

func FindCustomerByCompanyAndName(company *Company, name string) (*Customer, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Customer{}, "select * from customers where company_id = ? and name = ?", company.ID, name)
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

func CreateCustomer(
	company *Company,
	name string,
	kana string,
	birthday int64,
	gender Gender,
	postal_code1 string,
	postal_code2 string,
	address string,
	phone_number string,
	join_date int64,
	email_address string,
	can_mail bool,
	can_email bool,
	can_call bool,
	description string,
) (*Customer, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	customer, err := CreateCustomerInTx(
		company,
		name,
		kana,
		birthday,
		gender,
		postal_code1,
		postal_code2,
		address,
		phone_number,
		join_date,
		email_address,
		can_mail,
		can_email,
		can_call,
		description,
		tx,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return customer, nil
}

func CreateCustomerInTx(
	company *Company,
	name string,
	kana string,
	birthday int64,
	gender Gender,
	postal_code1 string,
	postal_code2 string,
	address string,
	phone_number string,
	join_date int64,
	email_address string,
	can_mail bool,
	can_email bool,
	can_call bool,
	description string,
	tx *gorp.Transaction,
) (*Customer, error) {
	if len(postal_code1) > 3 {
		return nil, fmt.Errorf("invalid length of postal_code1")
	}
	if len(postal_code2) > 4 {
		return nil, fmt.Errorf("invalid length of postal_code2")
	}
	customer := &Customer{
		CompanyID:    company.ID,
		Name:         name,
		Kana:         kana,
		Birthday:     birthday,
		Gender:       gender,
		PostalCode1:  postal_code1,
		PostalCode2:  postal_code2,
		Address:      address,
		PhoneNumber:  phone_number,
		JoinDate:     join_date,
		EmailAddress: email_address,
		CanMail:      can_mail,
		CanEmail:     can_email,
		CanCall:      can_call,
		Description:  description,
		CreatedAt:    time.Now().Unix(),
		IsValid:      true,
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
