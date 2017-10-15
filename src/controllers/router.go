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
	"SetCompanyImageLink": {
		"requireAuthorization": true,
		"delegate":             setCompanyImageLink{},
	},
	"SetCompanyPassword": {
		"requireAuthorization": true,
		"delegate":             setCompanyPassword{},
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
	"SelectStaffs": {
		"requireAuthorization": true,
		"delegate":             selectStaffs{},
	},
	"CreateStaff": {
		"requireAuthorization": true,
		"delegate":             createStaff{},
	},
	"UpdateStaff": {
		"requireAuthorization": true,
		"delegate":             updateStaff{},
	},
	"DeleteStaff": {
		"requireAuthorization": true,
		"delegate":             deleteStaff{},
	},
	"SelectCustomers": {
		"requireAuthorization": true,
		"delegate":             selectCustomers{},
	},
	"CreateCustomer": {
		"requireAuthorization": true,
		"delegate":             createCustomer{},
	},
	"UpdateCustomer": {
		"requireAuthorization": true,
		"delegate":             updateCustomer{},
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
