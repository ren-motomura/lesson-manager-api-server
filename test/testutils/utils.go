package testutils

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
)

type fakeReadCloser struct {
	io.Reader
}

func (fakeReadCloser) Close() error {
	return nil
}

func BuildRequest(funcName string, body []byte, sessionID string) (r *http.Request) {
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
