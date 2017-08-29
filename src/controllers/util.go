package controllers

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

func writeErrorResponseWithLog(err error, r *procesures.ParsedRequest, rw http.ResponseWriter, responseStatus int, errorType pb.ErrorType, errorMessage string) {
	ctx := appengine.NewContext(r.Origin)
	log.Errorf(ctx, "%s. error: %v", errorMessage, err)
	writeErrorResponse(rw, responseStatus, errorType, errorMessage)
}

func writeErrorResponse(rw http.ResponseWriter, responseStatus int, errorType pb.ErrorType, errorMessage string) {
	res, _ := proto.Marshal(&pb.ErrorResponse{ // エラーは発生しないはず
		ErrorType: errorType,
		Message:   errorMessage,
	})
	rw.WriteHeader(responseStatus)
	rw.Write(res)
}
