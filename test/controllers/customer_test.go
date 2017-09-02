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

func TestCreateCustomer(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	_, session := testutils.CreateCompanyAndSession()

	reqParam := &pb.CreateCustomerRequest{
		Name:        "sample customer",
		Description: "",
	}
	reqBin, _ := proto.Marshal(reqParam)
	req := testutils.BuildRequest("CreateCustomer", reqBin, session.ID)
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 200 {
		t.Fatalf("status: %d", frw.status)
	}
	res := &pb.CreateCustomerResponse{}
	err = proto.Unmarshal(frw.body, res)
	if err != nil {
		t.Fatal(err)
	}

	_, err = models.FindCustomer(int(res.Customer.Id), false, nil)
	if err != nil {
		t.Fatal(err)
	}

	if res.Customer.Card != nil {
		t.Fatal("card must be nil")
	}
}

func TestCreateCustomerWithCard(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	_, session := testutils.CreateCompanyAndSession()

	reqParam := &pb.CreateCustomerRequest{
		Name:        "sample customer",
		Description: "",
		Card: &pb.Card{
			Id:     "sample card",
			Credit: 10000,
		},
	}
	reqBin, _ := proto.Marshal(reqParam)
	req := testutils.BuildRequest("CreateCustomer", reqBin, session.ID)
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 200 {
		t.Fatalf("status: %d", frw.status)
	}
	res := &pb.CreateCustomerResponse{}
	err = proto.Unmarshal(frw.body, res)
	if err != nil {
		t.Fatal(err)
	}

	_, err = models.FindCustomer(int(res.Customer.Id), false, nil)
	if err != nil {
		t.Fatal(err)
	}

	if res.Customer.Card == nil {
		t.Fatal("card must not be nil")
	}

	_, err = models.FindCard(res.Customer.Card.Id, false, nil)
	if err != nil {
		t.Fatal(err)
	}
}
