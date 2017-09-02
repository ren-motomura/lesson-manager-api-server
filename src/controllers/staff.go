package controllers

import (
	"net/http"

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

	staff, err := models.CreateStaff(param.Name, r.Company)
	if err != nil {
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	res, _ := proto.Marshal(&pb.CreateStaffResponse{ // エラーは発生しないはず
		Staff: &pb.Staff{
			Id:        int32(staff.ID),
			Name:      staff.Name,
			CreatedAt: staff.CreatedAt.Unix(),
		},
	})
	rw.WriteHeader(200)
	rw.Write(res)
}
