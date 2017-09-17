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
		writeErrorResponse(rw, 409, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
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

type setCompanyImageLink struct {
}

func (setCompanyImageLink) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.SetCompanyImageLinkRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 409, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	r.Company.ImageLink = param.ImageLink
	err = r.Company.Update()
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 409, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	res, _ := proto.Marshal(&pb.SetCompanyImageLinkResponse{ // エラーは発生しないはず
		Company: &pb.Company{
			Id:           int32(r.Company.ID),
			Name:         r.Company.Name,
			EmailAddress: r.Company.EmailAddress,
			CreatedAt:    r.Company.CreatedAt.Unix(),
			ImageLInk:    r.Company.ImageLink,
		},
	})
	rw.WriteHeader(200)
	rw.Write(res)
}

type setCompanyPassword struct {
}

func (setCompanyPassword) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.SetCompanyPasswordRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 409, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	r.Company.SetPassword(param.Password)
	err = r.Company.Update()
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 409, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	res, _ := proto.Marshal(&pb.SetCompanyImageLinkResponse{ // エラーは発生しないはず
		Company: &pb.Company{
			Id:           int32(r.Company.ID),
			Name:         r.Company.Name,
			EmailAddress: r.Company.EmailAddress,
			CreatedAt:    r.Company.CreatedAt.Unix(),
			ImageLInk:    r.Company.ImageLink,
		},
	})
	rw.WriteHeader(200)
	rw.Write(res)
}
