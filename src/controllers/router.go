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
	"SelectStudios": {
		"requireAuthorization": true,
		"delegate":             selectStudios{},
	},
	"CreateStudio": {
		"requireAuthorization": true,
		"delegate":             createStudio{},
	},
	"UpdateStudio": {
		"requireAuthorization": true,
		"delegate":             updateStudio{},
	},
	"DeleteStudio": {
		"requireAuthorization": true,
		"delegate":             deleteStudio{},
	},
	"CreateStaff": {
		"requireAuthorization": true,
		"delegate":             createStaff{},
	},
	"DeleteStaff": {
		"requireAuthorization": true,
		"delegate":             deleteStaff{},
	},
	"CreateCustomer": {
		"requireAuthorization": true,
		"delegate":             createCustomer{},
	},
	"DeleteCustomer": {
		"requireAuthorization": true,
		"delegate":             deleteCustomer{},
	},
	"SetCardOnCustomer": {
		"requireAuthorization": true,
		"delegate":             setCardOnCustomer{},
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
