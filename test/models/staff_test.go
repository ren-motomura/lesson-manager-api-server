package models_test

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestStaff(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

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

	{
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

	{
		staffs[0].Delete()
		selectedStaffs, err := models.SelectStaffsByCompany(company)
		if err != nil {
			t.Fatal(err)
		}

		if len(staffs)-len(selectedStaffs) != 1 {
			t.Fatalf("len(staffs): %d, len(selectedStaffs): %d", len(staffs), len(selectedStaffs))
		}
	}
}
