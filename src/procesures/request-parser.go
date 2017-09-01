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
	Company  *models.Company
	Origin   *http.Request
}

const FuncNameHeaderKey = "X-Lessonmanager-Funcname"

func (req *ParsedRequest) IsAuthorized() bool {
	return req.Company != nil
}

func ParseRequest(r *http.Request) (*ParsedRequest, error) {
	if len(r.Header[FuncNameHeaderKey]) == 0 {
		return nil, errors.New("func name not found")
	}
	funcName := r.Header[FuncNameHeaderKey][0]
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("fail to read body")
	}

	company, _ := Authorize(r) // エラーは無視して大丈夫

	pr := ParsedRequest{
		FuncName: funcName,
		Data:     data,
		Company:  company,
		Origin:   r,
	}
	return &pr, nil
}
