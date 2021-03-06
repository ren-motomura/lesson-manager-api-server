package controllers_test

import (
	"testing"

	"github.com/ren-motomura/lesson-manager-api-server/src/errs"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/controllers"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestSelectStaffs(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	req := testutils.BuildRequest("SelectStaffs", []byte{}, session.ID)
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
		res := &pb.SelectStaffsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Staffs) != 0 {
			t.Fatal()
		}
	}

	staff, err := models.CreateStaff("sample staff", "", company)
	if err != nil {
		t.Fatal(err)
	}

	{
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SelectStaffsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Staffs) != 1 {
			t.Fatal()
		}
		if res.Staffs[0].Name != staff.Name {
			t.Fatal()
		}
	}
}

func TestCreateStaff(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	studioName := "sample studio"
	_, session := testutils.CreateCompanyAndSession()

	reqParam := &pb.CreateStaffRequest{
		Name:      studioName,
		ImageLink: "",
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

func TestUpdateStaff(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	updatedImageLink := "http://example.com/image"

	{
		staff, err := models.CreateStaff("sample staff", "", company)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.UpdateStaffRequest{
			Staff: &pb.Staff{
				Id:        int32(staff.ID),
				Name:      staff.Name,
				ImageLink: updatedImageLink,
			},
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("UpdateStaff", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}

		result, err := models.FindStaff(staff.ID, false, nil)
		if result.ImageLink != updatedImageLink {
			t.Fatal()
		}
	}

	{ // other company staff
		otherCompany, err := models.CreateCompany("sample company2", "sample2@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}

		otherCompanyStaff, err := models.CreateStaff("sample staff2", "", otherCompany)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.UpdateStaffRequest{
			Staff: &pb.Staff{
				Id:        int32(otherCompanyStaff.ID),
				Name:      otherCompanyStaff.Name,
				ImageLink: updatedImageLink,
			},
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("UpdateStaff", reqBin, session.ID)
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

func TestDeleteStaff(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	{
		staff, err := models.CreateStaff("sample staff", "", company)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.DeleteStaffRequest{
			Id: int32(staff.ID),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("DeleteStaff", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}

		_, err = models.FindStaff(staff.ID, false, nil)
		if err != errs.ErrNotFound {
			t.Fatal("not deleted")
		}
	}

	{ // other company staff
		otherCompany, err := models.CreateCompany("sample company2", "sample2@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}

		otherCompanyStaff, err := models.CreateStaff("sample staff2", "", otherCompany)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.DeleteStaffRequest{
			Id: int32(otherCompanyStaff.ID),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("DeleteStaff", reqBin, session.ID)
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
