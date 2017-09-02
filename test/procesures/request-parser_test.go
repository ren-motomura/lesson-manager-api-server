package procesures_test

import (
	"reflect"
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestParseRequest(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	funcName := "hoge"
	data := []byte{1, 2, 3, 4, 5}

	req := testutils.BuildRequest(funcName, data, "")

	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if pr.FuncName != funcName {
		t.Fatal()
	}

	if !reflect.DeepEqual(pr.Data, data) {
		t.Fatal()
	}
}
