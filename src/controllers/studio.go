package controllers

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

type createStudio struct {
}

func (createStudio) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.CreateStudioRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	studio, err := models.CreateStudio(param.Name, r.Company)
	if err != nil {
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	res, _ := proto.Marshal(&pb.CreateStudioResponse{ // エラーは発生しないはず
		Studio: &pb.Studio{
			Id:        int32(studio.ID),
			Name:      studio.Name,
			CreatedAt: studio.CreatedAt.Unix(),
		},
	})
	rw.WriteHeader(200)
	rw.Write(res)
}
