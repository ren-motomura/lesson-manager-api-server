package models_test

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestCard(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	name := "サードダンススクール"
	emailAddress := "third@example.com"
	password := fmt.Sprintf("%x", sha256.Sum256([]byte("password")))
	company, err := models.CreateCompany(name, emailAddress, password)
	if err != nil {
		t.Fatal(err)
	}

	customer, err := models.CreateCustomer(
		"sample customer",
		"description",
		company,
	)
	if err != nil {
		t.Fatal(err)
	}

	card, err := models.CreateCard("sample card", customer, 10000)
	if err != nil {
		t.Fatal(err)
	}

	_, err = models.FindCard(card.ID, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = models.FindCardByCustomer(customer, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	{
		customer2, err := models.CreateCustomer(
			"sample customer2",
			"description",
			company,
		)
		if err != nil {
			t.Fatal(err)
		}

		_, err = models.CreateCard("sample card2", customer2, 10000)
		if err != nil {
			t.Fatal(err)
		}

		cards, err := models.SelectCardsByCustomerIds([]int{customer.ID, customer2.ID})
		if err != nil {
			t.Fatal(err)
		}
		if len(cards) != 2 {
			t.Fatal()
		}
	}

	card.Delete()

	_, err = models.FindCard(card.ID, false, nil)
	if err != errs.ErrNotFound {
		t.Fatal("fail to delete")
	}
}
