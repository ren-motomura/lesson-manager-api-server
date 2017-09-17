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
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	company, err := models.FindCompanyByEmailAddress(param.EmailAddress)
	if err != nil {
		if err == errs.ErrNotFound {
			writeErrorResponse(rw, 404, pb.ErrorType_NOT_FOUND, "")
			return
		}
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	if !company.ComparePassword(param.Password) {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_PASSWORD, "")
		return
	}

	session, err := models.CreateSession(company)
	if err != nil { // まれに生成した sessionId がすでに存在していてエラーになる可能性があるが...その場合は 500 で返す
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	setSessionToResponse(rw, session)
	res, _ := proto.Marshal(&pb.CreateSessionResponse{ // エラーは発生しないはず
		Company: &pb.Company{
			Id:           int32(company.ID),
			Name:         company.Name,
			EmailAddress: company.EmailAddress,
			CreatedAt:    company.CreatedAt.Unix(),
			ImageLInk:    company.ImageLink,
		},
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
