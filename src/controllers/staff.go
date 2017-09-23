package controllers

import (
	"net/http"

	"github.com/ren-motomura/lesson-manager-api-server/src/errs"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

type createStaff struct {
}

func (createStaff) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.CreateStaffRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	staff, err := models.CreateStaff(param.Name, param.ImageLink, r.Company)
	if err != nil {
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
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
