package procesures_test

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
)

type fakeReadCloser struct {
	io.Reader
}

func (fakeReadCloser) Close() error {
	return nil
}

func TestParseRequest(t *testing.T) {
	funcName := "hoge"
	data := []byte{1, 2, 3, 4, 5}

	req := buildRequest(funcName, data, "")

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

func buildRequest(funcName string, body []byte, sessionID string) (r *http.Request) {
	r = new(http.Request)
	r.Header = make(map[string][]string)
	r.Header.Add(procesures.FuncNameHeaderKey, funcName)
	r.Body = &fakeReadCloser{bytes.NewBuffer(body)}
	r.AddCookie(&http.Cookie{
		Name:  procesures.SessionCookieName,
		Value: sessionID,
	})
	return r
}
