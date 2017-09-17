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

func TestCreateCompany(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	reqParam := &pb.CreateCompanyRequest{
		Name:         "サンプル",
		EmailAddress: "sample@example.com",
		Password:     "password",
	}
	reqBin, _ := proto.Marshal(reqParam)
	req := testutils.BuildRequest("CreateCompany", reqBin, "")
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 200 {
		t.Fatalf("status: %d", frw.status)
	}
	res := &pb.CreateCompanyResponse{}
	err = proto.Unmarshal(frw.body, res)
	if err != nil {
		t.Fatal(err)
	}
	resCompany := res.Company
	if resCompany.Name != reqParam.Name {
		t.Fatal()
	}
	if resCompany.EmailAddress != reqParam.EmailAddress {
		t.Fatal()
	}
}

func TestSetCompanyImageLink(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	_, session := testutils.CreateCompanyAndSession()

	reqParam := &pb.SetCompanyImageLinkRequest{
		ImageLink: "http://example.com/image",
	}
	reqBin, _ := proto.Marshal(reqParam)
	req := testutils.BuildRequest("SetCompanyImageLink", reqBin, session.ID)
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 200 {
		t.Fatalf("status: %d", frw.status)
	}
	res := &pb.SetCompanyImageLinkResponse{}
	err = proto.Unmarshal(frw.body, res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Company.ImageLInk != reqParam.ImageLink {
		t.Fatal()
	}
}

func TestSetCompanyPassword(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	reqParam := &pb.SetCompanyPasswordRequest{
		Password: "__newpassword",
	}
	reqBin, _ := proto.Marshal(reqParam)
	req := testutils.BuildRequest("SetCompanyPassword", reqBin, session.ID)
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 200 {
		t.Fatalf("status: %d", frw.status)
	}
	res := &pb.SetCompanyImageLinkResponse{}
	err = proto.Unmarshal(frw.body, res)
	if err != nil {
		t.Fatal(err)
	}

	company2, err := models.FindCompany(company.ID)
	if err != nil {
		t.Fatal(err)
	}

	if company.Password == company2.Password {
		t.Fatal()
	}
}
