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
}

func SelectLessonsByCompanyAndTakenAtRange(company *Company, takenAtFrom time.Time, takenAtTo time.Time) ([]*Lesson, error) {
	db, err := Db()
	if err != nil {
		return nil, err
	}

	rows, err := db.Select(Lesson{}, "select * from lessons where company_id = ? and taken_at >= ? and taken_at < ?", company.ID, takenAtFrom, takenAtTo)
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

	rows, err := db.Select(Lesson{}, "select * from lessons where staff_id = ? and taken_at >= ? and taken_at < ?", staff.ID, takenAtFrom, takenAtTo)
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

	rows, err := db.Select(Lesson{}, "select * from lessons where customer_id = ? and taken_at >= ? and taken_at < ?", customer.ID, takenAtFrom, takenAtTo)
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

	rows, err := db.Select(Lesson{}, "select * from lessons where studio_id = ? and taken_at >= ? and taken_at < ?", studio.ID, takenAtFrom, takenAtTo)
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
	}
	err := tx.Insert(lesson)
	if err != nil {
		return nil, err
	}

	return lesson, nil
}
