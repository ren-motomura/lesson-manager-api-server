package models_test

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestStudio(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	name := "サードダンススクール"
	emailAddress := "third@example.com"
	password := fmt.Sprintf("%x", sha256.Sum256([]byte("password")))
	company, err := models.CreateCompany(name, emailAddress, password)
	if err != nil {
		t.Fatal(err)
	}

	studios := make([]*models.Studio, 3)
	for i := 0; i < len(studios); i++ {
		studio, err := models.CreateStudio("sample studio"+strconv.Itoa(i), company)
		if err != nil {
			t.Fatal(err)
		}
		studios[i] = studio
	}

	selectedStudios, err := models.SelectStudiosByCompany(company)
	if err != nil {
		t.Fatal(err)
	}

	if len(studios) != len(selectedStudios) {
		t.Fatal(err)
	}

	for i, s := range selectedStudios {
		if studios[i].Name != s.Name {
			t.Fatal("invalid name")
		}
	}
}
