package controllers_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/controllers"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestCreateUser(t *testing.T) {
	reqParam := &pb.CreateUserRequest{
		Name:         "サンプル太郎",
		EmailAddress: "sample@expamle.com",
		Password:     "password",
	}
	reqBin, _ := proto.Marshal(reqParam)
	req := testutils.BuildRequest("CreateUser", reqBin, "")
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 200 {
		t.Fatalf("status: %d", frw.status)
	}
	res := &pb.CreateUserResponse{}
	err = proto.Unmarshal(frw.body, res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Name != reqParam.Name {
		t.Fatal()
	}
	if res.EmailAddress != reqParam.EmailAddress {
		t.Fatal()
	}
}
