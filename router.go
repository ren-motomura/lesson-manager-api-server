package router

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func init() {
	http.Handle("/", goji.DefaultMux)

	goji.Get("/", handler)
}

func handler(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "It's lesson manager!")
}
