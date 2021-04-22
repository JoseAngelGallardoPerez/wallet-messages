package rpc

import (
	"github.com/Confialink/wallet-messages/internal/dao"
	"github.com/Confialink/wallet-messages/internal/model"
	"context"
	"fmt"
	"net/http"

	"github.com/Confialink/wallet-messages/internal/config"
	pb "github.com/Confialink/wallet-messages/rpc/messages"
)

type PbServer struct {
	dao    *dao.Message
	config *config.Config
}

func NewPbServer(dao *dao.Message, config *config.Config) *PbServer {
	return &PbServer{dao, config}
}

func (s *PbServer) Start() {
	twirpHandler := pb.NewMessageSenderServer(s, nil)
	mux := http.NewServeMux()
	mux.Handle(pb.MessageSenderPathPrefix, twirpHandler)
	go http.ListenAndServe(fmt.Sprintf(":%s", s.config.ProtoBufPort), mux)
}

func (s *PbServer) SendMessage(ctx context.Context, req *pb.SendMessageReq) (*pb.MessageResp, error) {
	// check if there is the same unread message
	if req.DoNotDuplicateIfExists {
		message, err := s.dao.FindUnreadBySubjectAndSender(req.GetRecipientId(), req.GetSenderId(), req.GetSubject())
		if err != nil {
			return nil, err
		}

		if message.ID > 0 {
			return s.createResponse(message), nil
		}
	}
	message := model.Message{
		MessagePrivate: model.MessagePrivate{
			SenderId:            &req.SenderId,
			IsRecipientIncoming: true,
		},
		MessagePublic: model.MessagePublic{
			RecipientId:     &req.RecipientId,
			Subject:         &req.Subject,
			Message:         req.Message,
			DeleteAfterRead: &req.DeleteAfterRead,
		},
	}

	createdMessage, err := s.dao.Create(&message)

	if err != nil {
		return nil, err
	}

	return s.createResponse(createdMessage), nil
}

func (s *PbServer) createResponse(message *model.Message) *pb.MessageResp {
	return &pb.MessageResp{
		Id:          message.ID,
		RecipientId: *message.RecipientId,
		SenderId:    *message.SenderId,
		Subject:     *message.Subject,
		Message:     message.Message,
	}
}
