package models

import (
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"

	"github.com/go-gorp/gorp"
)

type Card struct {
	ID         string
	CustomerID int
	Credit     int
}

func registerCard(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Card{}, "cards").SetKeys(false, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("CustomerID").Rename("customer_id")
	t.ColMap("Credit").Rename("credit")
}

func (self *Card) Delete() error {
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

func (self *Card) DeleteInTx(tx *gorp.Transaction) error {
	_, err := tx.Delete(self)
	if err != nil {
		return err
	}

	return nil
}

func FindCard(id string, forUpdate bool, tx *gorp.Transaction) (*Card, error) {
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

	rows, err := selector.Select(Card{}, "select * from cards where id = ? "+forUpStatement, id)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}
	card := rows[0].(*Card)
	return card, nil
}

func CreateCard(id string, customer *Customer, credit int) (*Card, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	card, err := CreateCardInTx(id, customer, credit, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return card, nil
}

func CreateCardInTx(id string, customer *Customer, credit int, tx *gorp.Transaction) (*Card, error) {
	card := &Card{
		ID:         id,
		CustomerID: customer.ID,
		Credit:     credit,
	}
	err := tx.Insert(card)
	if err != nil {
		return nil, err
	}

	return card, nil
}
