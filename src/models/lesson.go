package models

import (
	"time"

	"github.com/go-gorp/gorp"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
)

type PaymentType int8

const (
	PaymentTypeCash PaymentType = iota
	PaymentTypeCard
)

type Lesson struct {
	ID          int
	CompanyID   int
	StudioID    int
	StaffID     int
	CustomerID  int
	Fee         int
	PaymentType PaymentType
	TakenAt     time.Time
	IsValid     bool
}

func registerLesson(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Lesson{}, "lessons").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("CompanyID").Rename("company_id")
	t.ColMap("StudioID").Rename("studio_id")
	t.ColMap("StaffID").Rename("staff_id")
	t.ColMap("CustomerID").Rename("customer_id")
	t.ColMap("Fee").Rename("fee")
	t.ColMap("PaymentType").Rename("payment_type")
	t.ColMap("TakenAt").Rename("taken_at")
	t.ColMap("IsValid").Rename("is_valid")
}

func (self *Lesson) Delete() error {
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

func (self *Lesson) DeleteInTx(tx *gorp.Transaction) error {
	self.IsValid = false
	_, err := tx.Update(self)
	if err != nil {
		return err
	}

	return nil
}

func FindLesson(id int, forUpdate bool, tx *gorp.Transaction) (*Lesson, error) {
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

	rows, err := selector.Select(Lesson{}, "select * from lessons where id = ? and is_valid = ? "+forUpStatement, id, true)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, errs.ErrNotFound
	}
	lesson := rows[0].(*Lesson)
	return lesson, nil
}

func SelectLessonsByCompanyAndTakenAtRange(company *Company, takenAtFrom time.Time, takenAtTo time.Time) ([]*Lesson, error) {
	if company == nil {
		return []*Lesson{}, nil
	}

	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Lesson{}, "select * from lessons where is_valid = ? and company_id = ? and taken_at >= ? and taken_at < ?", true, company.ID, takenAtFrom, takenAtTo)
	if err != nil {
		return nil, err
	}
	lessons := make([]*Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = row.(*Lesson)
	}
	return lessons, nil
}

func SelectLessonsByStaffAndTakenAtRange(staff *Staff, takenAtFrom time.Time, takenAtTo time.Time) ([]*Lesson, error) {
	if staff == nil {
		return []*Lesson{}, nil
	}

	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Lesson{}, "select * from lessons where is_valid = ? and staff_id = ? and taken_at >= ? and taken_at < ?", true, staff.ID, takenAtFrom, takenAtTo)
	if err != nil {
		return nil, err
	}
	lessons := make([]*Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = row.(*Lesson)
	}
	return lessons, nil
}

func SelectLessonsByCustomerAndTakenAtRange(customer *Customer, takenAtFrom time.Time, takenAtTo time.Time) ([]*Lesson, error) {
	if customer == nil {
		return []*Lesson{}, nil
	}

	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Lesson{}, "select * from lessons where is_valid = ? and customer_id = ? and taken_at >= ? and taken_at < ?", true, customer.ID, takenAtFrom, takenAtTo)
	if err != nil {
		return nil, err
	}
	lessons := make([]*Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = row.(*Lesson)
	}
	return lessons, nil
}

func SelectLessonsByStudioAndTakenAtRange(studio *Studio, takenAtFrom time.Time, takenAtTo time.Time) ([]*Lesson, error) {
	if studio == nil {
		return []*Lesson{}, nil
	}

	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Lesson{}, "select * from lessons where is_valid = ? and studio_id = ? and taken_at >= ? and taken_at < ?", true, studio.ID, takenAtFrom, takenAtTo)
	if err != nil {
		return nil, err
	}
	lessons := make([]*Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = row.(*Lesson)
	}
	return lessons, nil
}

func SelectLessonsByStudioAndStaffAndTakenAtRange(studio *Studio, staff *Staff, takenAtFrom time.Time, takenAtTo time.Time) ([]*Lesson, error) {
	if studio == nil || staff == nil {
		return []*Lesson{}, nil
	}

	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Lesson{}, "select * from lessons where is_valid = ? and studio_id = ? and staff_id = ? and taken_at >= ? and taken_at < ?", true, studio.ID, staff.ID, takenAtFrom, takenAtTo)
	if err != nil {
		return nil, err
	}
	lessons := make([]*Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = row.(*Lesson)
	}
	return lessons, nil
}

func SelectLessonsByStudioAndCustomerAndTakenAtRange(studio *Studio, customer *Customer, takenAtFrom time.Time, takenAtTo time.Time) ([]*Lesson, error) {
	if studio == nil || customer == nil {
		return []*Lesson{}, nil
	}

	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Lesson{}, "select * from lessons where is_valid = ? and studio_id = ? and customer_id = ? and taken_at >= ? and taken_at < ?", true, studio.ID, customer.ID, takenAtFrom, takenAtTo)
	if err != nil {
		return nil, err
	}
	lessons := make([]*Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = row.(*Lesson)
	}
	return lessons, nil
}

func SelectLessonsByStaffAndCustomerAndTakenAtRange(staff *Staff, customer *Customer, takenAtFrom time.Time, takenAtTo time.Time) ([]*Lesson, error) {
	if staff == nil || customer == nil {
		return []*Lesson{}, nil
	}

	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Lesson{}, "select * from lessons where is_valid = ? and staff_id = ? and customer_id = ? and taken_at >= ? and taken_at < ?", true, staff.ID, customer.ID, takenAtFrom, takenAtTo)
	if err != nil {
		return nil, err
	}
	lessons := make([]*Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = row.(*Lesson)
	}
	return lessons, nil
}

func SelectLessonsByStudioAndStaffAndCustomerAndTakenAtRange(studio *Studio, staff *Staff, customer *Customer, takenAtFrom time.Time, takenAtTo time.Time) ([]*Lesson, error) {
	if studio == nil || staff == nil || customer == nil {
		return []*Lesson{}, nil
	}

	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Lesson{}, "select * from lessons where is_valid = ? and studio_id = ? and staff_id = ? and customer_id = ? and taken_at >= ? and taken_at < ?", true, studio.ID, staff.ID, customer.ID, takenAtFrom, takenAtTo)
	if err != nil {
		return nil, err
	}
	lessons := make([]*Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = row.(*Lesson)
	}
	return lessons, nil
}

func CreateLesson(company *Company, studio *Studio, staff *Staff, customer *Customer, fee int, paymentType PaymentType, takenAt time.Time) (*Lesson, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	lesson, err := CreateLessonInTx(company, studio, staff, customer, fee, paymentType, takenAt, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return lesson, nil
}

func CreateLessonInTx(company *Company, studio *Studio, staff *Staff, customer *Customer, fee int, paymentType PaymentType, takenAt time.Time, tx *gorp.Transaction) (*Lesson, error) {
	lesson := &Lesson{
		CompanyID:   company.ID,
		StudioID:    studio.ID,
		StaffID:     staff.ID,
		CustomerID:  customer.ID,
		Fee:         fee,
		PaymentType: paymentType,
		TakenAt:     takenAt.UTC(),
		IsValid:     true,
	}
	err := tx.Insert(lesson)
	if err != nil {
		return nil, err
	}

	return lesson, nil
}
