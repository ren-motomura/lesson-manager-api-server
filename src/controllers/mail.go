package controllers

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/ren-motomura/lesson-manager-api-server/src/procesures"
	pb "github.com/ren-motomura/lesson-manager-api-server/src/protobufs"

	"google.golang.org/appengine"
	"google.golang.org/appengine/mail"
)

type sendMail struct {
}

func (sendMail) Execute(rw http.ResponseWriter, r *procesures.ParsedRequest) {
	param := &pb.SendEmailRequest{}
	err := proto.Unmarshal(r.Data, param)
	if err != nil {
		writeErrorResponse(rw, 400, pb.ErrorType_INVALID_REQUEST_FORMAT, "")
		return
	}

	if len(param.ToAddresses) == 0 {
		writeErrorResponse(rw, 400, pb.ErrorType_BAD_REQUEST, "")
		return
	}

	ctx := appengine.NewContext(r.Origin)
	msg := &mail.Message{
		Sender:  "LessonManager <noreply@third-being-175805.appspotmail.com>",
		To:      param.ToAddresses,
		Subject: param.Subject,
		Body:    param.Body,
	}
	attachments := make([]mail.Attachment, 0, len(param.Attachments))
	for _, at := range param.Attachments {
		attachments = append(attachments, mail.Attachment{
			Name: at.Name,
			Data: at.Data,
		})
	}
	msg.Attachments = attachments
	err = mail.Send(ctx, msg)
	if err != nil {
		writeErrorResponseWithLog(err, r, rw, 500, pb.ErrorType_INTERNAL_SERVER_ERROR, "")
		return
	}

	res, _ := proto.Marshal(&pb.SendEmailResponse{
		Success: true,
	})
	rw.WriteHeader(200)
	rw.Write(res)
}
