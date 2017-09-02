package controllers_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/controllers"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestCreateStudio(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	studioName := "sample studio"
	_, session := testutils.CreateCompanyAndSession()

	reqParam := &pb.CreateStudioRequest{
		Name: studioName,
	}
	reqBin, _ := proto.Marshal(reqParam)
	req := testutils.BuildRequest("CreateStudio", reqBin, session.ID)
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 200 {
		t.Fatalf("status: %d", frw.status)
	}
	res := &pb.CreateStudioResponse{}
	err = proto.Unmarshal(frw.body, res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Studio.Name != studioName {
		t.Fatal("unexpected name")
	}
}

func TestDeleteStudio(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	{
		studio, err := models.CreateStudio("sample studio", company)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.DeleteStudioRequest{
			Id: int32(studio.ID),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("DeleteStudio", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}

		_, err = models.FindStudio(studio.ID, false, nil)
		if err != errs.ErrNotFound {
			t.Fatal("not deleted")
		}
	}

	{ // other company studio
		otherCompany, err := models.CreateCompany("sample company2", "sample2@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}

		otherCompanyStudio, err := models.CreateStudio("sample studio2", otherCompany)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.DeleteStudioRequest{
			Id: int32(otherCompanyStudio.ID),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("DeleteStudio", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 403 {
			t.Fatalf("status: %d", frw.status)
		}
	}

}
