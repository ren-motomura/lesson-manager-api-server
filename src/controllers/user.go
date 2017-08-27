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
		rw.WriteHeader(409)
		return
	}
	if err != errs.ErrNotFound { // 謎のエラー
		ctx := appengine.NewContext(r.Origin)
		log.Errorf(ctx, "find user error: %v", err)
		rw.WriteHeader(500)
		return
	}

	// ユーザーとセッションを同一トランザクションで作成する
	db, err := models.Db()
	if err != nil {
		ctx := appengine.NewContext(r.Origin)
		log.Errorf(ctx, "db connection error: %v", err)
		rw.WriteHeader(500)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ctx := appengine.NewContext(r.Origin)
		log.Errorf(ctx, "start transaction error: %v", err)
		rw.WriteHeader(500)
		return
	}

	user, err := models.CreateUserInTx(param.Name, param.EmailAddress, param.Password, tx)
	if err != nil {
		tx.Rollback()
		ctx := appengine.NewContext(r.Origin)
		log.Errorf(ctx, "create user error: %v", err)
		rw.WriteHeader(500)
		return
	}

	session, err := models.CreateSessionInTx(user, tx)
	if err != nil {
		tx.Rollback()
		ctx := appengine.NewContext(r.Origin)
		log.Errorf(ctx, "create session error: %v", err)
		rw.WriteHeader(500)
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
