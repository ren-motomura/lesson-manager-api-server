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
		studio, err := models.CreateStudio("sample studio"+strconv.Itoa(i), "sample address", "00-0000-0000", company, "")
		if err != nil {
			t.Fatal(err)
		}
		studios[i] = studio
	}

	{
		_, err = models.FindStudioByCompanyAndName(company, studios[0].Name)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
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

	{
		studios[0].Delete()
		selectedStudios, err := models.SelectStudiosByCompany(company)
		if err != nil {
			t.Fatal(err)
		}

		if len(studios)-len(selectedStudios) != 1 {
			t.Fatal("fail to delete")
		}
	}

	{
		s := studios[1]
		s.Address = "updated address"
		s.PhoneNumber = "updated number"
		s.ImageLink = "http://example.com/image"
		err = s.Update()
		if err != nil {
			t.Fatal(err)
		}

		res, err := models.FindStudio(s.ID, false, nil)
		if err != nil {
			t.Fatal(err)
		}

		if res.Address != s.Address || res.PhoneNumber != s.PhoneNumber || res.ImageLink != s.ImageLink {
			t.Fatal()
		}
	}
}
