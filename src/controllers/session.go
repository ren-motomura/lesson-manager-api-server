package controllers

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

type createSession struct {
}

func (createSession) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.CreateSessionRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		rw.WriteHeader(400)
		return
	}

	user, err := models.FindUserByEmailAddress(param.EmailAddress)
	if err != nil {
		if err == errs.ErrNotFound {
			rw.WriteHeader(404)
			return
		}
		rw.WriteHeader(500)
		return
	}
	session, err := models.CreateSession(user)
	if err != nil { // まれに生成した sessionId がすでに存在していてエラーになる可能性があるが...その場合は 500 で返す
		rw.WriteHeader(500)
		return
	}

	cookie := &http.Cookie{
		Name:  procesures.SessionCookieName,
		Value: session.ID,
	}
	http.SetCookie(rw, cookie)
	rw.WriteHeader(200)
	res, err := proto.Marshal(&pb.CreateSessionResponse{
		Success: true,
	})
	if err != nil {
		ctx := appengine.NewContext(r.Origin)
		log.Errorf(ctx, "protobuf marshal error: %v", err)
		rw.WriteHeader(500)
		return
	}
	rw.WriteHeader(200)
	rw.Write(res)
}
