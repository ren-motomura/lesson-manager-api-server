package models_test

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestLesson(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	name := "サードダンススクール"
	emailAddress := "third@example.com"
	password := fmt.Sprintf("%x", sha256.Sum256([]byte("password")))
	company, err := models.CreateCompany(name, emailAddress, password)
	if err != nil {
		t.Fatal(err)
	}

	studios := make([]*models.Studio, 2)
	for i := 0; i < len(studios); i++ {
		studio, err := models.CreateStudio(
			"sample studio"+strconv.Itoa(i),
			"sample address",
			"00-0000-0000",
			company,
			"",
		)
		if err != nil {
			t.Fatal(err)
		}
		studios[i] = studio
	}

	staffs := make([]*models.Staff, 3)
	for i := 0; i < len(staffs); i++ {
		staff, err := models.CreateStaff(
			"sample staff"+strconv.Itoa(i),
			"",
			company,
		)
		if err != nil {
			t.Fatal(err)
		}
		staffs[i] = staff
	}

	customers := make([]*models.Customer, 4)
	for i := 0; i < len(customers); i++ {
		customer, err := models.CreateCustomer(
			"sample customer"+strconv.Itoa(i),
			"description",
			company,
		)
		if err != nil {
			t.Fatal(err)
		}
		customers[i] = customer
	}

	now := time.Now()
	lessons := make([]*models.Lesson, 0, len(studios)*len(staffs)*len(customers))
	for _, studio := range studios {
		for _, staff := range staffs {
			for _, customer := range customers {
				lesson, err := models.CreateLesson(
					company,
					studio,
					staff,
					customer,
					6000,
					models.PaymentTypeCard,
					now,
				)
				if err != nil {
					t.Fatal(err)
				}
				lessons = append(lessons, lesson)
			}
		}
	}

	{ // 単体取得
		lesson, err := models.FindLesson(lessons[0].ID, false, nil)
		if err != nil {
			t.Fatal(err)
		}
		if lesson.StaffID != lessons[0].StaffID {
			t.Fatal()
		}
	}

	{ // 時間が対象外
		selectedLessons, err := models.SelectLessonsByCompanyAndTakenAtRange(company, now.Add(time.Second), now.Add(2*time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != 0 {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}
	}

	{ // 全部取得
		selectedLessons, err := models.SelectLessonsByCompanyAndTakenAtRange(company, now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != len(studios)*len(staffs)*len(customers) {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}

		selectedLessons, _ = models.SelectLessonsByCompanyAndTakenAtRange(nil, now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}
	}

	{ // スタジオで取得
		selectedLessons, err := models.SelectLessonsByStudioAndTakenAtRange(studios[0], now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != len(staffs)*len(customers) {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}

		selectedLessons, _ = models.SelectLessonsByStudioAndTakenAtRange(nil, now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}
	}

	{ // スタッフで取得
		selectedLessons, err := models.SelectLessonsByStaffAndTakenAtRange(staffs[0], now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != len(studios)*len(customers) {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}

		selectedLessons, _ = models.SelectLessonsByStaffAndTakenAtRange(nil, now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}
	}

	{ // 顧客で取得
		selectedLessons, err := models.SelectLessonsByCustomerAndTakenAtRange(customers[0], now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != len(studios)*len(staffs) {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}

		selectedLessons, _ = models.SelectLessonsByCustomerAndTakenAtRange(nil, now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}
	}

	{ // スタジオとスタッフで取得
		selectedLessons, err := models.SelectLessonsByStudioAndStaffAndTakenAtRange(studios[0], staffs[0], now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != len(customers) {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}

		selectedLessons, _ = models.SelectLessonsByStudioAndStaffAndTakenAtRange(nil, staffs[0], now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}

		selectedLessons, _ = models.SelectLessonsByStudioAndStaffAndTakenAtRange(studios[0], nil, now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}
	}

	{ // スタジオと顧客で取得
		selectedLessons, err := models.SelectLessonsByStudioAndCustomerAndTakenAtRange(studios[0], customers[0], now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != len(staffs) {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}

		selectedLessons, _ = models.SelectLessonsByStudioAndCustomerAndTakenAtRange(nil, customers[0], now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}

		selectedLessons, _ = models.SelectLessonsByStudioAndCustomerAndTakenAtRange(studios[0], nil, now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}
	}

	{ // スタッフと顧客で取得
		selectedLessons, err := models.SelectLessonsByStaffAndCustomerAndTakenAtRange(staffs[0], customers[0], now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != len(studios) {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}

		selectedLessons, _ = models.SelectLessonsByStaffAndCustomerAndTakenAtRange(nil, customers[0], now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}

		selectedLessons, _ = models.SelectLessonsByStaffAndCustomerAndTakenAtRange(staffs[0], nil, now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}
	}

	{ // スタジオとスタッフと顧客で取得
		selectedLessons, err := models.SelectLessonsByStudioAndStaffAndCustomerAndTakenAtRange(studios[0], staffs[0], customers[0], now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != 1 {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}

		selectedLesson := selectedLessons[0]

		if selectedLesson.StudioID != studios[0].ID {
			t.Fatal()
		}

		if selectedLesson.StaffID != staffs[0].ID {
			t.Fatal()
		}

		if selectedLesson.CustomerID != customers[0].ID {
			t.Fatal()
		}

		selectedLessons, _ = models.SelectLessonsByStudioAndStaffAndCustomerAndTakenAtRange(nil, staffs[0], customers[0], now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}

		selectedLessons, _ = models.SelectLessonsByStudioAndStaffAndCustomerAndTakenAtRange(studios[0], nil, customers[0], now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}

		selectedLessons, _ = models.SelectLessonsByStudioAndStaffAndCustomerAndTakenAtRange(studios[0], staffs[0], nil, now.Add(-time.Second), now.Add(time.Second))
		if len(selectedLessons) != 0 {
			t.Fatal()
		}
	}

	{ // 削除のテスト
		lessons[0].Delete()
		selectedLessons, err := models.SelectLessonsByCompanyAndTakenAtRange(company, now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != len(studios)*len(staffs)*len(customers)-1 {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}
	}
}
