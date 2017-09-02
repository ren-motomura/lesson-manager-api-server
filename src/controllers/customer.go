package controllers

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/errs"
	"github.com/ren-motomura/lesson-manager-api-server/src/models"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"
)

type createCustomer struct {
}

func (createCustomer) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.CreateCustomerRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	if param.Card != nil {
		_, err = models.FindCard(param.Card.Id, false, nil)
		if err != errs.ErrNotFound {
			writeErrorResponse(rw, 409, pb.ErrorType_ALREADY_EXIST, "")
			return
		}
	}

	db, err := models.Db()
	if err != nil {
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	customer, err := models.CreateCustomerInTx(param.Name, param.Description, r.Company, tx)
	if err != nil {
		tx.Rollback()
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	var card *models.Card
	if param.Card != nil {
		card, err = models.CreateCardInTx(param.Card.Id, customer, int(param.Card.Credit), tx)
		if err != nil {
			tx.Rollback()
			writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		writeErrorResponse(rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	var pbCard *pb.Card
	if card != nil {
		pbCard = &pb.Card{
			Id:     card.ID,
			Credit: int32(card.Credit),
		}
	}
	res, _ := proto.Marshal(&pb.CreateCustomerResponse{ // エラーは発生しないはず
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
