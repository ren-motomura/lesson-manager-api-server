package models_test

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
)

func TestUser(t *testing.T) {
	name := "田中太郎"
	emailAddress := "taro.tanaka@example.com"
	password := fmt.Sprintf("%x", sha256.Sum256([]byte("password")))
	user, err := models.CreateUser(name, emailAddress, password)
	if err != nil {
		t.Fatal(err)
	}

	user, err = models.FindUserByEmailAddress(emailAddress)
	if err != nil {
		t.Fatal(err)
	}
	if user.Name != name {
		t.Fatal(
			"Invalid name was returned! expected: ",
			name,
			", actual: ",
			user.Name,
		)
	}

	user, err = models.FindUser(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if user.Name != name {
		t.Fatal(
			"Invalid name was returned! expected: ",
			name,
			", actual: ",
			user.Name,
		)
	}
}
