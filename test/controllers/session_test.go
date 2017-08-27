package controllers_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/controllers"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestCreateSession(t *testing.T) {
	user, err := models.CreateUser("sample太郎", "sample@example.com", "password")
	if err != nil {
		t.Fatal(err)
	}

	reqParam := &pb.CreateSessionRequest{
		EmailAddress: user.EmailAddress,
		Password:     user.Password,
	}
	reqBin, err := proto.Marshal(reqParam)
	req := testutils.BuildRequest("CreateSession", reqBin, "")
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 200 {
		t.Fatal()
	}
	res := &pb.CreateSessionResponse{}
	err = proto.Unmarshal(frw.body, res)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Success {
		t.Fatal()
	}
}
