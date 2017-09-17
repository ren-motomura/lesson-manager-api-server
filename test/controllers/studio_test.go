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

func TestSelectStudios(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	req := testutils.BuildRequest("SelectStudios", []byte{}, session.ID)
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	{
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SelectStudiosResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Studios) != 0 {
			t.Fatal()
		}
	}

	studio, err := models.CreateStudio("sample studio", "sample address", "00-0000-0000", company, "")
	if err != nil {
		t.Fatal(err)
	}

	{
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SelectStudiosResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Studios) != 1 {
			t.Fatal()
		}
		if res.Studios[0].Name != studio.Name {
			t.Fatal()
		}
	}
}

func TestCreateStudio(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	studioName := "sample studio"
	_, session := testutils.CreateCompanyAndSession()

	reqParam := &pb.CreateStudioRequest{
		Name:        studioName,
		Address:     "sample address",
		PhoneNumber: "00-0000-0000",
		ImageLink:   "",
	}
	reqBin, _ := proto.Marshal(reqParam)

	req := testutils.BuildRequest("CreateStudio", reqBin, session.ID)
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	{
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

	{
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 409 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.ErrorResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}
		if res.ErrorType != pb.ErrorType_ALREADY_EXIST {
			t.Fatal()
		}
	}
}

func TestUpdateStudio(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	updatedAddress := "updated address"
	updatedPhoneNumber := "11-1111-1111"
	updatedImageLink := "http://example.com/image"

	{
		studio, err := models.CreateStudio("sample studio", "sample address", "00-0000-0000", company, "")
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.UpdateStudioRequest{
			Studio: &pb.Studio{
				Id:          int32(studio.ID),
				Name:        studio.Name,
				Address:     updatedAddress,
				PhoneNumber: updatedPhoneNumber,
				ImageLink:   updatedImageLink,
				CreatedAt:   studio.CreatedAt.Unix(),
			},
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("UpdateStudio", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}

		result, err := models.FindStudio(studio.ID, false, nil)
		if result.Address != updatedAddress || result.PhoneNumber != updatedPhoneNumber || result.ImageLink != updatedImageLink {
			t.Fatal()
		}
	}

	{ // other company studio
		otherCompany, err := models.CreateCompany("sample company2", "sample2@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}

		otherCompanyStudio, err := models.CreateStudio("sample studio2", "sample address", "00-0000-0000", otherCompany, "")
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.UpdateStudioRequest{
			Studio: &pb.Studio{
				Id:          int32(otherCompanyStudio.ID),
				Name:        otherCompanyStudio.Name,
				Address:     updatedAddress,
				PhoneNumber: updatedPhoneNumber,
				ImageLink:   updatedImageLink,
				CreatedAt:   otherCompanyStudio.CreatedAt.Unix(),
			},
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("UpdateStudio", reqBin, session.ID)
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

func TestDeleteStudio(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	{
		studio, err := models.CreateStudio("sample studio", "sample address", "00-0000-0000", company, "")
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

		otherCompanyStudio, err := models.CreateStudio("sample studio2", "sample address", "00-0000-0000", otherCompany, "")
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
