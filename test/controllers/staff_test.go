package controllers_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/controllers"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestCreateStaff(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	studioName := "sample studio"
	_, session := testutils.CreateCompanyAndSession()

	reqParam := &pb.CreateStaffRequest{
		Name: studioName,
	}
	reqBin, _ := proto.Marshal(reqParam)
	req := testutils.BuildRequest("CreateStaff", reqBin, session.ID)
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 200 {
		t.Fatalf("status: %d", frw.status)
	}
	res := &pb.CreateStaffResponse{}
	err = proto.Unmarshal(frw.body, res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Staff.Name != studioName {
		t.Fatal("unexpected name")
	}
}
