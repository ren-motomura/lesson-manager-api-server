package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// CompaniesController ...
type CompaniesController struct {
}

// SetHandler ...
func (c *CompaniesController) SetHandler(r *mux.Router, prefix string) {
	r.HandleFunc(prefix, c.getCompanies)
}

func (c *CompaniesController) getCompanies(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "get companies")
}
