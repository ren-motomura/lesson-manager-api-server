package procesures_test

import (
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestAuthorize(t *testing.T) {
	company, err := models.CreateCompany("test", "test@example.com", "password")
	if err != nil {
		t.Fatal(err)
	}

	session, err := models.CreateSession(company)
	if err != nil {
		t.Fatal(err)
	}

	req := testutils.BuildRequest("testFunc", []byte{}, session.ID)

	authorizedCompany, err := procesures.Authorize(req)
	if err != nil {
		t.Fatal("authorize failed")
	}

	if authorizedCompany.ID != company.ID {
		t.Fatal("something wrong")
	}
}
