package models_test

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
)

func TestStaff(t *testing.T) {
	name := "サードダンススクール"
	emailAddress := "third@example.com"
	password := fmt.Sprintf("%x", sha256.Sum256([]byte("password")))
	company, err := models.CreateCompany(name, emailAddress, password)
	if err != nil {
		t.Fatal(err)
	}

	staffs := make([]*models.Staff, 3)
	for i := 0; i < len(staffs); i++ {
		staff, err := models.CreateStaff("sample staff"+strconv.Itoa(i), company)
		if err != nil {
			t.Fatal(err)
		}
		staffs[i] = staff
	}

	selectedStaffs, err := models.SelectStaffsByCompany(company)
	if err != nil {
		t.Fatal(err)
	}

	if len(staffs) != len(selectedStaffs) {
		t.Fatal(err)
	}

	for i, s := range selectedStaffs {
		if staffs[i].Name != s.Name {
			t.Fatal("invalid name")
		}
	}
}
