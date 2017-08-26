package procesures_test

import (
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
)

func TestAuthorize(t *testing.T) {
	user, err := models.CreateUser("test user", "test@example.com", "password")
	if err != nil {
		t.Fatal(err)
	}

	session, err := models.CreateSession(user)
	if err != nil {
		t.Fatal(err)
	}

	req := buildRequest("testFunc", []byte{}, session.ID)

	authorizedUser, err := procesures.Authorize(req)
	if err != nil {
		t.Fatal("authorize failed")
	}

	if authorizedUser.ID != user.ID {
		t.Fatal("something wrong")
	}
}
