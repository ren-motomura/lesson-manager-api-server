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
			company,
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
			"card"+strconv.Itoa(i),
			0,
		)
		if err != nil {
			t.Fatal(err)
		}
		customers[i] = customer
	}

	now := time.Now()
	for _, studio := range studios {
		for _, staff := range staffs {
			for _, customer := range customers {
				_, err := models.CreateLesson(
					company,
					studio,
					staff,
					customer,
					6000,
					now,
				)
				if err != nil {
					t.Fatal(err)
				}
			}
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
	}

	{ // スタジオで取得
		selectedLessons, err := models.SelectLessonsByStudioAndTakenAtRange(studios[0], now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != len(staffs)*len(customers) {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
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
	}

	{ // 顧客で取得
		selectedLessons, err := models.SelectLessonsByCustomerAndTakenAtRange(customers[0], now.Add(-time.Second), now.Add(time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if len(selectedLessons) != len(studios)*len(staffs) {
			t.Fatalf("Invalid count: %d", len(selectedLessons))
		}
	}
}
