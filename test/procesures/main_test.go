package procesures_test

import (
	"os"
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestMain(m *testing.M) {
	testutils.PreProcess()
	code := m.Run()
	testutils.PostProcess()
	os.Exit(code)
}
