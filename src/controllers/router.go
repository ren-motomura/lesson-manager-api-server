package controllers

import (
	"net/http"

	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
)

var funcMap = map[string]map[string]interface{}{
	"CreateCompany": {
		"requireAuthorization": false,
		"delegate":             createCompany{},
	},
	"CreateSession": {
		"requireAuthorization": false,
		"delegate":             createSession{},
	},
	"CreateStudio": {
		"requireAuthorization": true,
		"delegate":             createStudio{},
	},
	"CreateStaff": {
		"requireAuthorization": true,
		"delegate":             createStaff{},
	},
}

type executer interface {
	Execute(http.ResponseWriter, *procesures.ParsedRequest)
}

var withoutAuthorizeFuncs = []string{
	"CreateSession",
}

func Route(rw http.ResponseWriter, req *procesures.ParsedRequest) {
	funcData := funcMap[req.FuncName]
	if funcData == nil {
		rw.WriteHeader(404)
		return
	}
	if !req.IsAuthorized() && funcData["requireAuthorization"].(bool) {
		rw.WriteHeader(401)
		return
	}
	funcData["delegate"].(executer).Execute(rw, req)
}
