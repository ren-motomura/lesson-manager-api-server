package procesures

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
)

type ParsedRequest struct {
	FuncName string
	Data     []byte
	User     *models.User
}

const FuncNameHeaderKey = "X-Lessonmanager-Funcname"

func ParseRequest(r *http.Request) (*ParsedRequest, error) {
	if len(r.Header[FuncNameHeaderKey]) == 0 {
		return nil, errors.New("func name not found")
	}
	funcName := r.Header[FuncNameHeaderKey][0]
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("fail to read body")
	}

	user, _ := Authorize(r) // エラーは無視して大丈夫

	pr := ParsedRequest{
		FuncName: funcName,
		Data:     data,
		User:     user,
	}
	return &pr, nil
}
