package controllers

import (
	"net/http"

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
		res, _ := proto.Marshal(&pb.ErrorResponse{ // エラーは発生しないはず
			ErrorType: pb.ErrorType_INVALID_REQUEST_FORMAT,
			Message:   "",
		})
		rw.WriteHeader(400)
		rw.Write(res)
		return
	}

	user, err := models.FindUserByEmailAddress(param.EmailAddress)
	if err != nil {
		if err == errs.ErrNotFound {
			res, _ := proto.Marshal(&pb.ErrorResponse{ // エラーは発生しないはず
				ErrorType: pb.ErrorType_USER_NOT_FOUND,
				Message:   "",
			})
			rw.WriteHeader(404)
			rw.Write(res)
			return
		}
		res, _ := proto.Marshal(&pb.ErrorResponse{ // エラーは発生しないはず
			ErrorType: pb.ErrorType_INTERNAL_SERVER_ERROR,
			Message:   "",
		})
		rw.WriteHeader(500)
		rw.Write(res)
		return
	}

	if !user.ComparePassword(param.Password) {
		res, _ := proto.Marshal(&pb.ErrorResponse{ // エラーは発生しないはず
			ErrorType: pb.ErrorType_INVALID_PASSWORD,
			Message:   "",
		})
		rw.WriteHeader(400)
		rw.Write(res)
		return
	}

	session, err := models.CreateSession(user)
	if err != nil { // まれに生成した sessionId がすでに存在していてエラーになる可能性があるが...その場合は 500 で返す
		res, _ := proto.Marshal(&pb.ErrorResponse{ // エラーは発生しないはず
			ErrorType: pb.ErrorType_INTERNAL_SERVER_ERROR,
			Message:   "",
		})
		rw.WriteHeader(500)
		rw.Write(res)
		return
	}

	setSessionToResponse(rw, session)
	res, _ := proto.Marshal(&pb.CreateSessionResponse{ // エラーは発生しないはず
		Success: true,
	})
	rw.WriteHeader(200)
	rw.Write(res)
}

func setSessionToResponse(rw http.ResponseWriter, session *models.Session) {
	cookie := &http.Cookie{
		Name:  procesures.SessionCookieName,
		Value: session.ID,
	}
	http.SetCookie(rw, cookie)
}
