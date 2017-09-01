package controllers

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

type createCompany struct {
}

func (createCompany) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.CreateCompanyRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		rw.WriteHeader(400)
		return
	}

	_, err = models.FindCompanyByEmailAddress(param.EmailAddress)
	if err == nil { // Company が存在したことになる
		writeErrorResponse(rw, 409, pb.ErrorType_ALREADY_EXIST, "")
		return
	}
	if err != errs.ErrNotFound { // 謎のエラー
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "fail to find company")
		return
	}

	// Company と Session を同一トランザクションで作成する
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

	company, err := models.CreateCompanyInTx(param.Name, param.EmailAddress, param.Password, tx)
	if err != nil {
		tx.Rollback()
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "fail to create company")
		return
	}

	session, err := models.CreateSessionInTx(company, tx)
	if err != nil {
		tx.Rollback()
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "fail to create session")
		return
	}

	tx.Commit()

	setSessionToResponse(rw, session)
	res, _ := proto.Marshal(&pb.CreateCompanyResponse{ // エラーは発生しないはず
		Id:           int64(company.ID),
		Name:         company.Name,
		EmailAddress: company.EmailAddress,
		CreatedAt:    company.CreatedAt.Unix(),
	})
	rw.WriteHeader(200)
	rw.Write(res)
}
