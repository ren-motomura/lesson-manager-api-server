package controllers_test

import (
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/controllers"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestRouteWithNotExistFunc(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	req := testutils.BuildRequest("notexist", []byte{}, "")
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 404 {
		t.Fatal()
	}
}
