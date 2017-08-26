package models_test

import (
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
)

func TestDb(t *testing.T) {
	_, err := models.Db()
	if err != nil {
		t.Fatal(err)
	}
}
