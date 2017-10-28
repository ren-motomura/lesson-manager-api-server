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

func TestSelectCustomers(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	customer1, err := models.CreateCustomer("customer1", "desc", company)
	if err != nil {
		t.Fatal(err)
	}
	customer2, err := models.CreateCustomer("customer2", "desc", company)
	if err != nil {
		t.Fatal(err)
	}

	// customer1 にだけカードを登録しておく
	card, err := models.CreateCard("card1", customer1, 100)
	if err != nil {
		t.Fatal(err)
	}

	req := testutils.BuildRequest("SelectCustomers", []byte{}, session.ID)
	pr, err := procesures.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	frw := fakeResponseWriter{}
	controllers.Route(&frw, pr)

	if frw.status != 200 {
		t.Fatalf("status: %d", frw.status)
	}
	res := &pb.SelectCustomersResponse{}
	err = proto.Unmarshal(frw.body, res)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Customers) != 2 {
		t.Fatal()
	}

	if res.Customers[0].Name != customer1.Name {
		t.Fatal()
	}

	if res.Customers[1].Name != customer2.Name {
		t.Fatal()
	}

	if res.Customers[0].Card == nil {
		t.Fatal()
	}

	if res.Customers[0].Card.Id != card.ID {
		t.Fatal()
	}

	if int(res.Customers[0].Card.Credit) != card.Credit {
		t.Fatal()
	}

	if res.Customers[1].Card != nil {
		t.Fatal()
	}
}

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
	{ // 正常系
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

	{ // 名前の重複エラー
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

		if res.ErrorType != pb.ErrorType_DUPLICATE_NAME_EXIST {
			t.Fatal()
		}
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

func TestUpdateCustomer(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	updatedName := "update name"
	updatedDescription := "updated description"
	{
		customer, err := models.CreateCustomer("sample customer", "desc", company)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.UpdateCustomerRequest{
			Customer: &pb.Customer{
				Id:          int32(customer.ID),
				Name:        updatedName,
				Description: updatedDescription,
			},
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("UpdateCustomer", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}

		newCustomer, err := models.FindCustomer(customer.ID, false, nil)
		if err != nil {
			t.Fatal(err)
		}

		if newCustomer.Name != updatedName {
			t.Fatal()
		}

		if newCustomer.Description != updatedDescription {
			t.Fatal()
		}
	}

	{ // 名前の重複エラー
		customer, err := models.CreateCustomer("sample customer2", "desc", company)
		reqParam := &pb.UpdateCustomerRequest{
			Customer: &pb.Customer{
				Id:          int32(customer.ID),
				Name:        updatedName,
				Description: updatedDescription,
			},
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("UpdateCustomer", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
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

		if res.ErrorType != pb.ErrorType_DUPLICATE_NAME_EXIST {
			t.Fatal()
		}
	}
}

func TestDeleteCustomer(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	{
		customer, err := models.CreateCustomer("sample customer", "desc", company)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.DeleteCustomerRequest{
			Id: int32(customer.ID),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("DeleteCustomer", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}

		_, err = models.FindCustomer(customer.ID, false, nil)
		if err != errs.ErrNotFound {
			t.Fatal("not deleted")
		}
	}

	{ // other company customer
		otherCompany, err := models.CreateCompany("sample company2", "sample2@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}

		otherCompanyCustomer, err := models.CreateCustomer("sample customer2", "desc", otherCompany)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.DeleteStudioRequest{
			Id: int32(otherCompanyCustomer.ID),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("DeleteCustomer", reqBin, session.ID)
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

func TestSetCardOnCustomer(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	customer, err := models.CreateCustomer("sample customer", "desc", company)
	if err != nil {
		t.Fatal(err)
	}

	var card *models.Card
	{
		reqParam := &pb.SetCardOnCustomerRequest{
			CustomerId: int32(customer.ID),
			Card: &pb.Card{
				Id:     "sample card",
				Credit: 100000,
			},
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SetCardOnCustomer", reqBin, session.ID)
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

		card, err = models.FindCard(res.Customer.Card.Id, false, nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		reqParam := &pb.SetCardOnCustomerRequest{
			CustomerId: int32(customer.ID),
			Card: &pb.Card{
				Id:     card.ID,
				Credit: 100000,
			},
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SetCardOnCustomer", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 409 {
			t.Fatalf("status: %d", frw.status)
		}
	}
}

func TestAddCredit(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	amount := 100
	{
		customer, err := models.CreateCustomer("sample customer", "desc", company)
		if err != nil {
			t.Fatal(err)
		}

		card, err := models.CreateCard("sample card", customer, 0)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.AddCreditRequest{
			CustomerId: int32(customer.ID),
			Amount:     int32(amount),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("AddCredit", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.AddCreditResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		card, err = models.FindCard(res.Customer.Card.Id, false, nil)
		if err != nil {
			t.Fatal(err)
		}
		if card.Credit != amount {
			t.Fatal()
		}
	}

	{
		customer, err := models.CreateCustomer("sample customer2", "desc", company)
		if err != nil {
			t.Fatal(err)
		}

		reqParam := &pb.AddCreditRequest{
			CustomerId: int32(customer.ID),
			Amount:     int32(amount),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("AddCredit", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 400 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.ErrorResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if res.ErrorType != pb.ErrorType_CARD_NOT_REGISTERED {
			t.Fatal()
		}
	}
}
