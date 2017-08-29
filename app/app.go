package app

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"

	"github.com/ren-motomura/lesson-manager-api-server/src/controllers"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
	"github.com/urfave/negroni"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)

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

func mainHandler(rw http.ResponseWriter, r *http.Request) {
	req, err := procesures.ParseRequest(r)
	if err != nil { // ここでのエラーはリクエストの形式がおかしい場合のみ
		res, _ := proto.Marshal(&pb.ErrorResponse{ // エラーは発生しないはず
			ErrorType: pb.ErrorType_INVALID_REQUEST_FORMAT,
			Message:   "invalid format",
		})

		rw.WriteHeader(400)
		rw.Write(res)
		return
	}
	controllers.Route(rw, req)
}
