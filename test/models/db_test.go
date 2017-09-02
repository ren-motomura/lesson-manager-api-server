package models_test

import (
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestDb(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	_, err := models.Db()
	if err != nil {
		t.Fatal(err)
	}
}
