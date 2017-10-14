package controllers

import (
	"net/http"
	"time"

	"github.com/ren-motomura/lesson-manager-api-server/src/errs"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

type selectStaffs struct {
}

func (selectStaffs) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	staffs, err := models.SelectStaffsByCompany(r.Company)
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	var pbStaffs []*pb.Staff
	for _, s := range staffs {
		pbStaffs = append(pbStaffs, &pb.Staff{
			Id:        int32(s.ID),
			Name:      s.Name,
			ImageLink: s.ImageLink,
			CreatedAt: s.CreatedAt.Unix(),
		})
	}
	res, _ := proto.Marshal(&pb.SelectStaffsResponse{
		Staffs: pbStaffs,
	})
	rw.WriteHeader(200)
	rw.Write(res)
}

type createStaff struct {
}

func (createStaff) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.CreateStaffRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	_, err = models.FindStaffByCompanyAndName(r.Company, param.Name)
	if err != errs.ErrNotFound {
		writeErrorResponse(rw, 409, pb.ErrorType_ALREADY_EXIST, "")
		return
	}

	staff, err := models.CreateStaff(param.Name, param.ImageLink, r.Company)
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	res, _ := proto.Marshal(&pb.CreateStaffResponse{ // エラーは発生しないはず
		Staff: &pb.Staff{
			Id:        int32(staff.ID),
			Name:      staff.Name,
			CreatedAt: staff.CreatedAt.Unix(),
			ImageLink: staff.ImageLink,
		},
	})
	rw.WriteHeader(200)
	rw.Write(res)
}

type updateStaff struct {
}

func (updateStaff) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	time.Sleep(time.Second * 3)
	param := &pb.UpdateStaffRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	staff, err := models.FindStaff(int(param.Staff.Id), false, nil)
	if err != nil {
		if err == errs.ErrNotFound {
			writeErrorResponse(rw, 404, pb.ErrorType_NOT_FOUND, "")
			return
		}
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	if staff.CompanyID != r.Company.ID {
		writeErrorResponse(rw, 403, pb.ErrorType_FORBIDDEN, "")
		return
	}

	// 名前は更新しない
	staff.ImageLink = param.Staff.ImageLink

	err = staff.Update()
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	res, _ := proto.Marshal(&pb.UpdateStaffResponse{ // エラーは発生しないはず
		Staff: &pb.Staff{
			Id:        int32(staff.ID),
			Name:      staff.Name,
			CreatedAt: staff.CreatedAt.Unix(),
			ImageLink: staff.ImageLink,
		},
	})
	rw.WriteHeader(200)
	rw.Write(res)
}

type deleteStaff struct {
}

func (deleteStaff) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.DeleteStaffRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	staff, err := models.FindStaff(int(param.Id), false, nil)
	if err != nil {
		if err == errs.ErrNotFound {
			// もともと存在しないのならAPI的には成功
			res, _ := proto.Marshal(&pb.DeleteStaffResponse{
				Success: true,
			})
			rw.WriteHeader(200)
			rw.Write(res)
			return
		}
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	if staff.CompanyID != r.Company.ID {
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

	staff, err = models.FindStaff(int(param.Id), true, tx)
	if err != nil {

		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	err = staff.DeleteInTx(tx)
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

	res, _ := proto.Marshal(&pb.DeleteStaffResponse{ // エラーは発生しないはず
		Success: true,
	})
	rw.WriteHeader(200)
	rw.Write(res)
}
