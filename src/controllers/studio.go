package controllers

import (
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

type selectStudios struct {
}

func (selectStudios) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	studios, err := models.SelectStudiosByCompany(r.Company)
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	var pbStudios []*pb.Studio
	for _, s := range studios {
		pbStudios = append(pbStudios, &pb.Studio{
			Id:          int32(s.ID),
			Name:        s.Name,
			Address:     s.Address,
			PhoneNumber: s.PhoneNumber,
			ImageLink:   s.ImageLink,
			CreatedAt:   s.CreatedAt.Unix(),
		})
	}
	res, _ := proto.Marshal(&pb.SelectStudiosResponse{
		Studios: pbStudios,
	})
	rw.WriteHeader(200)
	rw.Write(res)
}

type createStudio struct {
}

func (createStudio) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.CreateStudioRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	_, err = models.FindStudioByCompanyAndName(r.Company, param.Name)
	if err != errs.ErrNotFound {
		writeErrorResponse(rw, 409, pb.ErrorType_ALREADY_EXIST, "")
		return
	}

	studio, err := models.CreateStudio(param.Name, param.Address, param.PhoneNumber, r.Company, param.ImageLink)
	if err != nil {
		// 存在チェックの直後に insert されるとエラーになり得るが、このサービスの用途から考えてまず起こらないので無視
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	res, _ := proto.Marshal(&pb.CreateStudioResponse{ // エラーは発生しないはず
		Studio: &pb.Studio{
			Id:        int32(studio.ID),
			Name:      studio.Name,
			CreatedAt: studio.CreatedAt.Unix(),
			ImageLink: studio.ImageLink,
		},
	})
	rw.WriteHeader(200)
	rw.Write(res)
}

type updateStudio struct {
}

func (updateStudio) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	time.Sleep(time.Second * 3)
	param := &pb.UpdateStudioRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	studio, err := models.FindStudio(int(param.Studio.Id), false, nil)
	if err != nil {
		if err == errs.ErrNotFound {
			writeErrorResponse(rw, 404, pb.ErrorType_NOT_FOUND, "")
			return
		}
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	if studio.CompanyID != r.Company.ID {
		writeErrorResponse(rw, 403, pb.ErrorType_FORBIDDEN, "")
		return
	}

	// 名前は更新しない
	studio.Address = param.Studio.Address
	studio.PhoneNumber = param.Studio.PhoneNumber
	studio.ImageLink = param.Studio.ImageLink

	err = studio.Update()
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	res, _ := proto.Marshal(&pb.UpdateStudioResponse{ // エラーは発生しないはず
		Studio: &pb.Studio{
			Id:          int32(studio.ID),
			Name:        studio.Name,
			Address:     studio.Address,
			PhoneNumber: studio.PhoneNumber,
			CreatedAt:   studio.CreatedAt.Unix(),
			ImageLink:   studio.ImageLink,
		},
	})
	rw.WriteHeader(200)
	rw.Write(res)
}

type deleteStudio struct {
}

func (deleteStudio) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.DeleteStudioRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	studio, err := models.FindStudio(int(param.Id), false, nil)
	if err != nil {
		if err == errs.ErrNotFound {
			// もともと存在しないのならAPI的には成功
			res, _ := proto.Marshal(&pb.DeleteStudioResponse{
				Success: true,
			})
			rw.WriteHeader(200)
			rw.Write(res)
			return
		}
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	if studio.CompanyID != r.Company.ID {
		writeErrorResponse(rw, 403, pb.ErrorType_FORBIDDEN, "")
		return
	}

	db, err := models.Db()
	if err != nil {
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	studio, err = models.FindStudio(int(param.Id), true, tx)
	if err != nil {

		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	err = studio.DeleteInTx(tx)
	if err != nil {
		tx.Rollback()
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	res, _ := proto.Marshal(&pb.DeleteStudioResponse{ // エラーは発生しないはず
		Success: true,
	})
	rw.WriteHeader(200)
	rw.Write(res)
}
