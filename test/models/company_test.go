package models_test

import (
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestCompany(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	name := "サードダンススクール"
	emailAddress := "third@example.com"
	password := "password"
	company, err := models.CreateCompany(name, emailAddress, password)
	if err != nil {
		t.Fatal(err)
	}

	company, err = models.FindCompanyByEmailAddress(emailAddress)
	if err != nil {
		t.Fatal(err)
	}
	if company.Name != name {
		t.Fatal(
			"Invalid name was returned! expected: ",
			name,
			", actual: ",
			company.Name,
		)
	}

	company, err = models.FindCompany(company.ID)
	if err != nil {
		t.Fatal(err)
	}
	if company.Name != name {
		t.Fatal(
			"Invalid name was returned! expected: ",
			name,
			", actual: ",
			company.Name,
		)
	}
}
