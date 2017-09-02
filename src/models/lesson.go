package models

import (
	"time"

	"github.com/go-gorp/gorp"
)

type Lesson struct {
	ID         int
	CompanyID  int
	StudioID   int
	StaffID    int
	CustomerID int
	Fee        int
	TakenAt    time.Time
	IsValid    bool
}

func registerLesson(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(Lesson{}, "lessons").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("CompanyID").Rename("company_id")
	t.ColMap("StudioID").Rename("studio_id")
	t.ColMap("StaffID").Rename("staff_id")
	t.ColMap("CustomerID").Rename("customer_id")
	t.ColMap("Fee").Rename("fee")
	t.ColMap("TakenAt").Rename("taken_at")
	t.ColMap("IsValid").Rename("is_valid")
}

func (self *Lesson) Delete() error {
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

func SelectLessonsByCompanyAndTakenAtRange(company *Company, takenAtFrom time.Time, takenAtTo time.Time) ([]*Lesson, error) {
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

func CreateLesson(company *Company, studio *Studio, staff *Staff, customer *Customer, fee int, takenAt time.Time) (*Lesson, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	tx, _ := db.Begin()
	lesson, err := CreateLessonInTx(company, studio, staff, customer, fee, takenAt, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return lesson, nil
}

func CreateLessonInTx(company *Company, studio *Studio, staff *Staff, customer *Customer, fee int, takenAt time.Time, tx *gorp.Transaction) (*Lesson, error) {
	lesson := &Lesson{
		CompanyID:  company.ID,
		StudioID:   studio.ID,
		StaffID:    staff.ID,
		CustomerID: customer.ID,
		Fee:        fee,
		TakenAt:    takenAt,
		IsValid:    true,
	}
	err := tx.Insert(lesson)
	if err != nil {
		return nil, err
	}

	return lesson, nil
}