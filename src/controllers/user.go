package controllers

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

type createUser struct {
}

func (createUser) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.CreateUserRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		rw.WriteHeader(400)
		return
	}

	_, err = models.FindUserByEmailAddress(param.EmailAddress)
	if err == nil { // ユーザーが存在したことになる
		writeErrorResponse(rw, 409, pb.ErrorType_USER_ALREADY_EXIST, "")
		return
	}
	if err != errs.ErrNotFound { // 謎のエラー
		writeErrorResponseWithLog(err, r, rw, 409, pb.ErrorType_USER_ALREADY_EXIST, "fail to find user")
		return
	}

	// ユーザーとセッションを同一トランザクションで作成する
	db, err := models.Db()
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "fail to get db")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "fail to start transaction")
		return
	}

	user, err := models.CreateUserInTx(param.Name, param.EmailAddress, param.Password, tx)
	if err != nil {
		tx.Rollback()
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "fail to create user")
		return
	}

	session, err := models.CreateSessionInTx(user, tx)
	if err != nil {
		tx.Rollback()
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "fail to create session")
		return
	}

	tx.Commit()

	setSessionToResponse(rw, session)
	res, _ := proto.Marshal(&pb.CreateUserResponse{ // エラーは発生しないはず
		Id:           int64(user.ID),
		Name:         user.Name,
		EmailAddress: user.EmailAddress,
	})
	rw.WriteHeader(200)
	rw.Write(res)
}
