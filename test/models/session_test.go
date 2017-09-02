package models_test

import (
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestSession(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	sessionInserted, err := models.CreateSession(&models.Company{
		Name:         "hogehogename",
		EmailAddress: "hoge@example.com",
		Password:     "hogehogepassword",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = models.FindSession(sessionInserted.ID)
	if err != nil {
		t.Fatal(err)
	}
}
