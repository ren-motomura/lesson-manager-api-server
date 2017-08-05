package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(myMiddleware))
	n.UseHandler(r)

	http.Handle("/", n)
}

func myMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "middleware!")
	next(rw, r)
}

func handler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "It's lesson manager!!")
}
