package controllers

import (
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

type registerLesson struct {
}

func (registerLesson) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.RegisterLessonRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	studio, err := models.FindStudio(int(param.StudioId), false, nil)
	if err != nil {
		if err == errs.ErrNotFound {
			writeErrorResponse(rw, 404, pb.ErrorType_STUDIO_NOT_FOUND, "")
			return
		}
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}
	if studio.CompanyID != r.Company.ID {
		writeErrorResponse(rw, 403, pb.ErrorType_FORBIDDEN, "")
		return
	}

	staff, err := models.FindStaff(int(param.StaffId), false, nil)
	if err != nil {
		if err == errs.ErrNotFound {
			writeErrorResponse(rw, 404, pb.ErrorType_STAFF_NOT_FOUND, "")
			return
		}
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}
	if staff.CompanyID != r.Company.ID {
		writeErrorResponse(rw, 403, pb.ErrorType_FORBIDDEN, "")
		return
	}

	customer, err := models.FindCustomer(int(param.CustomerId), false, nil)
	if err != nil {
		if err == errs.ErrNotFound {
			writeErrorResponse(rw, 404, pb.ErrorType_CUSTOMER_NOT_FOUND, "")
			return
		}
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}
	if customer.CompanyID != r.Company.ID {
		writeErrorResponse(rw, 403, pb.ErrorType_FORBIDDEN, "")
		return
	}

	var card *models.Card
	if param.PaymentType == pb.PaymentType_ByCard { // カードの存在チェック
		card, err = models.FindCardByCustomer(customer, false, nil)
		if err == errs.ErrNotFound {
			writeErrorResponse(rw, 400, pb.ErrorType_CARD_NOT_REGISTERED, "")
			return
		}
		if err != nil {
			writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
			return
		}

		if card.Credit < int(param.Fee) {
			writeErrorResponse(rw, 400, pb.ErrorType_CREDIT_SHORTAGE, "")
			return
		}
	}

	db, err := models.Db()
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	if param.PaymentType == pb.PaymentType_ByCard { // カードの残高処理
		card, err = models.FindCardByCustomer(customer, true, tx)
		if err == errs.ErrNotFound {
			writeErrorResponse(rw, 400, pb.ErrorType_CARD_NOT_REGISTERED, "")
			return
		}
		if err != nil {
			writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
			return
		}

		if card.Credit < int(param.Fee) {
			writeErrorResponse(rw, 400, pb.ErrorType_CREDIT_SHORTAGE, "")
			return
		}

		card.Credit = card.Credit - int(param.Fee)
		_, err = tx.Update(card)
		if err != nil {
			tx.Rollback()
			writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
			return
		}
	}

	lesson, err := models.CreateLessonInTx(
		r.Company,
		studio,
		staff,
		customer,
		int(param.Fee),
		pbPaymentTYpeToPaymentType(param.PaymentType),
		time.Unix(param.TakenAt, 0),
		tx,
	)
	if err != nil {
		tx.Rollback()
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	// レスポンス用にカード情報をまとめる
	if param.PaymentType != pb.PaymentType_ByCard {
		card, _ = models.FindCardByCustomer(customer, false, nil)
	}
	var pbCard *pb.Card
	if card != nil {
		pbCard = &pb.Card{
			Id:     card.ID,
			Credit: int32(card.Credit),
		}
	}

	res, _ := proto.Marshal(&pb.RegisterLessonResponse{
		Lesson: &pb.Lesson{
			Id:          int32(lesson.ID),
			StudioId:    int32(lesson.StudioID),
			StaffId:     int32(lesson.StaffID),
			CustomerId:  int32(lesson.CustomerID),
			Fee:         int32(lesson.Fee),
			PaymentType: paymentTypeToPbPaymentType(lesson.PaymentType),
			TakenAt:     lesson.TakenAt.Unix(),
		},
		Customer: &pb.Customer{
			Id:          int32(customer.ID),
			Name:        customer.Name,
			Description: customer.Description,
			Card:        pbCard,
		},
	})
	rw.WriteHeader(200)
	rw.Write(res)
}
