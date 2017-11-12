package controllers_test

import (
	"strconv"
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

func TestSearchLessons(t *testing.T) {
	teardown := testutils.Setup(t)
	defer teardown(t)

	company, session := testutils.CreateCompanyAndSession()

	studios := make([]*models.Studio, 2)
	for i := 0; i < len(studios); i++ {
		studio, err := models.CreateStudio(
			"sample studio"+strconv.Itoa(i),
			"sample address",
			"00-0000-0000",
			company,
			"",
		)
		if err != nil {
			t.Fatal(err)
		}
		studios[i] = studio
	}

	staffs := make([]*models.Staff, 3)
	for i := 0; i < len(staffs); i++ {
		staff, err := models.CreateStaff(
			"sample staff"+strconv.Itoa(i),
			"",
			company,
		)
		if err != nil {
			t.Fatal(err)
		}
		staffs[i] = staff
	}

	customers := make([]*models.Customer, 4)
	for i := 0; i < len(customers); i++ {
		customer, err := models.CreateCustomer(
			"sample customer"+strconv.Itoa(i),
			"description",
			company,
		)
		if err != nil {
			t.Fatal(err)
		}
		customers[i] = customer
	}

	now := time.Now()
	lessons := make([]*models.Lesson, 0, len(studios)*len(staffs)*len(customers))
	for _, studio := range studios {
		for _, staff := range staffs {
			for _, customer := range customers {
				lesson, err := models.CreateLesson(
					company,
					studio,
					staff,
					customer,
					6000,
					models.PaymentTypeCard,
					now,
				)
				if err != nil {
					t.Fatal(err)
				}
				lessons = append(lessons, lesson)
			}
		}
	}

	{ // 全指定
		reqParam := &pb.SearchLessonsRequest{
			StudioId:    int32(studios[0].ID),
			StaffId:     int32(staffs[0].ID),
			CustomerId:  int32(customers[0].ID),
			TakenAtFrom: now.Add(-time.Minute).Unix(),
			TakenAtTo:   now.Add(time.Minute).Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SearchLessons", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SearchLessonsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Lessons) != 1 {
			t.Fatal()
		}
		if int(res.Lessons[0].StudioId) != studios[0].ID {
			t.Fatal()
		}
		if int(res.Lessons[0].StaffId) != staffs[0].ID {
			t.Fatal()
		}
		if int(res.Lessons[0].CustomerId) != customers[0].ID {
			t.Fatal()
		}
	}

	{ // 全指定
		reqParam := &pb.SearchLessonsRequest{
			StudioId:    int32(studios[0].ID),
			StaffId:     int32(staffs[0].ID),
			CustomerId:  int32(customers[0].ID),
			TakenAtFrom: now.Add(-time.Minute).Unix(),
			TakenAtTo:   now.Add(time.Minute).Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SearchLessons", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SearchLessonsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Lessons) != 1 {
			t.Fatal()
		}
		if int(res.Lessons[0].StudioId) != studios[0].ID {
			t.Fatal()
		}
		if int(res.Lessons[0].StaffId) != staffs[0].ID {
			t.Fatal()
		}
		if int(res.Lessons[0].CustomerId) != customers[0].ID {
			t.Fatal()
		}
	}

	{ // スタジオ、スタッフ指定
		reqParam := &pb.SearchLessonsRequest{
			StudioId:    int32(studios[0].ID),
			StaffId:     int32(staffs[0].ID),
			CustomerId:  -1,
			TakenAtFrom: now.Add(-time.Minute).Unix(),
			TakenAtTo:   now.Add(time.Minute).Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SearchLessons", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SearchLessonsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Lessons) != len(customers) {
			t.Fatal()
		}
		if int(res.Lessons[0].StudioId) != studios[0].ID {
			t.Fatal()
		}
		if int(res.Lessons[0].StaffId) != staffs[0].ID {
			t.Fatal()
		}
	}

	{ // スタジオ、顧客指定
		reqParam := &pb.SearchLessonsRequest{
			StudioId:    int32(studios[0].ID),
			StaffId:     -1,
			CustomerId:  int32(customers[0].ID),
			TakenAtFrom: now.Add(-time.Minute).Unix(),
			TakenAtTo:   now.Add(time.Minute).Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SearchLessons", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SearchLessonsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Lessons) != len(staffs) {
			t.Fatal()
		}
		if int(res.Lessons[0].StudioId) != studios[0].ID {
			t.Fatal()
		}
		if int(res.Lessons[0].CustomerId) != customers[0].ID {
			t.Fatal()
		}
	}

	{ // スタッフ、顧客指定
		reqParam := &pb.SearchLessonsRequest{
			StudioId:    -1,
			StaffId:     int32(staffs[0].ID),
			CustomerId:  int32(customers[0].ID),
			TakenAtFrom: now.Add(-time.Minute).Unix(),
			TakenAtTo:   now.Add(time.Minute).Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SearchLessons", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SearchLessonsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Lessons) != len(studios) {
			t.Fatal()
		}
		if int(res.Lessons[0].StaffId) != staffs[0].ID {
			t.Fatal()
		}
		if int(res.Lessons[0].CustomerId) != customers[0].ID {
			t.Fatal()
		}
	}

	{ // スタジオ指定
		reqParam := &pb.SearchLessonsRequest{
			StudioId:    int32(studios[0].ID),
			StaffId:     -1,
			CustomerId:  -1,
			TakenAtFrom: now.Add(-time.Minute).Unix(),
			TakenAtTo:   now.Add(time.Minute).Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SearchLessons", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SearchLessonsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Lessons) != len(staffs)*len(customers) {
			t.Fatal()
		}
		if int(res.Lessons[0].StudioId) != studios[0].ID {
			t.Fatal()
		}
	}

	{ // スタッフ指定
		reqParam := &pb.SearchLessonsRequest{
			StudioId:    -1,
			StaffId:     int32(staffs[0].ID),
			CustomerId:  -1,
			TakenAtFrom: now.Add(-time.Minute).Unix(),
			TakenAtTo:   now.Add(time.Minute).Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SearchLessons", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SearchLessonsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Lessons) != len(studios)*len(customers) {
			t.Fatal()
		}
		if int(res.Lessons[0].StaffId) != staffs[0].ID {
			t.Fatal()
		}
	}

	{ // 顧客指定
		reqParam := &pb.SearchLessonsRequest{
			StudioId:    -1,
			StaffId:     -1,
			CustomerId:  int32(customers[0].ID),
			TakenAtFrom: now.Add(-time.Minute).Unix(),
			TakenAtTo:   now.Add(time.Minute).Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SearchLessons", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SearchLessonsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Lessons) != len(studios)*len(staffs) {
			t.Fatal()
		}
		if int(res.Lessons[0].CustomerId) != customers[0].ID {
			t.Fatal()
		}
	}

	{ // 無指定
		reqParam := &pb.SearchLessonsRequest{
			StudioId:    -1,
			StaffId:     -1,
			CustomerId:  -1,
			TakenAtFrom: now.Add(-time.Minute).Unix(),
			TakenAtTo:   now.Add(time.Minute).Unix(),
		}
		reqBin, _ := proto.Marshal(reqParam)
		req := testutils.BuildRequest("SearchLessons", reqBin, session.ID)
		pr, err := procesures.ParseRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		frw := fakeResponseWriter{}
		controllers.Route(&frw, pr)

		if frw.status != 200 {
			t.Fatalf("status: %d", frw.status)
		}
		res := &pb.SearchLessonsResponse{}
		err = proto.Unmarshal(frw.body, res)
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Lessons) != len(studios)*len(staffs)*len(customers) {
			t.Fatal()
		}
	}
}
