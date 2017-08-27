package utils_test

import (
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/utils"
)

func TestIndex(t *testing.T) {
	vs := []string{"a", "b", "c", "d", "e"}
	if utils.Index(vs, "e") != 4 {
		t.Fatal()
	}
	if utils.Index(vs, "f") != -1 {
		t.Fatal()
	}
}

func TestInclude(t *testing.T) {
	vs := []string{"a", "b", "c", "d", "e"}
	if !utils.Include(vs, "e") {
		t.Fatal()
	}
	if utils.Include(vs, "f") {
		t.Fatal()
	}
}
