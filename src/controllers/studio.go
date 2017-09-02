package controllers

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

type createStudio struct {
}

func (createStudio) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.CreateStudioRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		rw.WriteHeader(400)
		return
	}

	//company, err := models.FindCompany(int(param.CompanyId))
}
