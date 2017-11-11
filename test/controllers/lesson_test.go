package controllers_test

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/controllers"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
	"github.com/ren-motomura/lesson-manager-api-server/test/testutils"
)

func TestRegisterLesson(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()
	studio, _ := models.CreateStudio("studio1", "", "", company, "")
	staff, _ := models.CreateStaff("staff1", "", company)
	customer, _ := models.CreateCustomer("customer1", "desc", company)

	otherCompany, _ := models.CreateCompany("other company", "", "")
	otherCompanyStudio, _ := models.CreateStudio("studio2", "", "", otherCompany, "")
	otherCompanyStaff, _ := models.CreateStaff("staff2", "", otherCompany)
	otherCompanyCustomer, _ := models.CreateCustomer("customer2", "desc", otherCompany)

	{ // 正常系
		reqParam := &pb.RegisterLessonRequest{
			StudioId:   int32(studio.ID),
			StaffId:    int32(staff.ID),
			CustomerId: int32(customer.ID),
			Fee:        5000,
			TakenAt:    time.Now().Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("RegisterLesson", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.RegisterLessonResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		_, err = models.FindLesson(int(res.Lesson.Id), false, nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	{ // スタジオがない場合
		reqParam := &pb.RegisterLessonRequest{
			StudioId:   int32(studio.ID) + 100,
			StaffId:    int32(staff.ID),
			CustomerId: int32(customer.ID),
			Fee:        5000,
			TakenAt:    time.Now().Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("RegisterLesson", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 404 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.ErrorResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if res.ErrorType != pb.ErrorType_STUDIO_NOT_FOUND {
			t.Fatal()
		}
	}

	{ // スタッフがいない場合
		reqParam := &pb.RegisterLessonRequest{
			StudioId:   int32(studio.ID),
			StaffId:    int32(staff.ID) + 100,
			CustomerId: int32(customer.ID),
			Fee:        5000,
			TakenAt:    time.Now().Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("RegisterLesson", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 404 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.ErrorResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if res.ErrorType != pb.ErrorType_STAFF_NOT_FOUND {
			t.Fatal()
		}
	}

	{ // 顧客がいない場合
		reqParam := &pb.RegisterLessonRequest{
			StudioId:   int32(studio.ID),
			StaffId:    int32(staff.ID),
			CustomerId: int32(customer.ID) + 100,
			Fee:        5000,
			TakenAt:    time.Now().Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("RegisterLesson", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 404 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.ErrorResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if res.ErrorType != pb.ErrorType_CUSTOMER_NOT_FOUND {
			t.Fatal()
		}
	}

	{ // スタジオが自分のものでない場合
		reqParam := &pb.RegisterLessonRequest{
			StudioId:   int32(otherCompanyStudio.ID),
			StaffId:    int32(staff.ID),
			CustomerId: int32(customer.ID),
			Fee:        5000,
			TakenAt:    time.Now().Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("RegisterLesson", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 403 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.ErrorResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if res.ErrorType != pb.ErrorType_FORBIDDEN {
			t.Fatal()
		}
	}

	{ // スタッフが自分のものでない場合
		reqParam := &pb.RegisterLessonRequest{
			StudioId:   int32(studio.ID),
			StaffId:    int32(otherCompanyStaff.ID),
			CustomerId: int32(customer.ID),
			Fee:        5000,
			TakenAt:    time.Now().Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("RegisterLesson", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 403 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.ErrorResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if res.ErrorType != pb.ErrorType_FORBIDDEN {
			t.Fatal()
		}
	}

	{ // 顧客が自分のものでない場合
		reqParam := &pb.RegisterLessonRequest{
			StudioId:   int32(studio.ID),
			StaffId:    int32(staff.ID),
			CustomerId: int32(otherCompanyCustomer.ID),
			Fee:        5000,
			TakenAt:    time.Now().Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("RegisterLesson", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 403 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.ErrorResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if res.ErrorType != pb.ErrorType_FORBIDDEN {
			t.Fatal()
		}
	}
}
