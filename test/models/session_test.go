package models_test

import (
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
)

func TestSession(t *testing.T) {
	sessionInserted, err := models.CreateSession(&models.User{
		ID:           1,
		Name:         "hogehogename",
		EmailAddress: "hoge@example.com",
		Password:     "hogehogepassword",
	})
	if err != nil {
		t.Fatal(err)
	}

	sessionFound, err := models.FindSession(sessionInserted.ID)
	if err != nil {
		t.Fatal(err)
	}

	if sessionFound.ID != sessionInserted.ID {
		t.Fatal(
			"Invalid session id was returned! expected: ",
			sessionInserted.ID,
			", actual: ",
			sessionFound.ID,
		)
	}
}
