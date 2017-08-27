package controllers_test

import (
	"bytes"
	"net/http"
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

type fakeResponseWriter struct {
	status int
	body   []byte
}

func (*fakeResponseWriter) Header() http.Header {
	return http.Header{}
}
func (frw *fakeResponseWriter) Write(data []byte) (int, error) {
	buf := bytes.NewBuffer([]byte{})
	i, err := buf.Write(data)
	if err != nil {
		return i, err
	}
	frw.body = buf.Bytes()
	return i, err
}
func (frw *fakeResponseWriter) WriteHeader(status int) {
	frw.status = status
}
