package models_test

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestCustomer(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	name := "サードダンススクール"
	emailAddress := "third@example.com"
	password := fmt.Sprintf("%x", sha256.Sum256([]byte("password")))
	company, err := models.CreateCompany(name, emailAddress, password)
	if err != nil {
		t.Fatal(err)
	}

	customers := make([]*models.Customer, 3)
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

	selectedCustomers, err := models.SelectCustomersByCompany(company)
	if err != nil {
		t.Fatal(err)
	}

	if len(customers) != len(selectedCustomers) {
		t.Fatal(err)
	}

	for i, s := range selectedCustomers {
		if customers[i].Name != s.Name {
			t.Fatal("invalid name")
		}
	}

	customer, err := models.FindCustomerByCompanyAndCardID(company, "card0")
	if err != nil {
		t.Fatal(err)
	}

	if customer.ID != customers[0].ID {
		t.Fatal("invalid ID")
	}

	customers[0].Delete()
	_, err = models.FindCustomerByCompanyAndCardID(company, "card0")
	if err != errs.ErrNotFound {
		t.Fatal("fail to delete")
	}
}
