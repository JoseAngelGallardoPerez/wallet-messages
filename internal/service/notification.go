package service

import (
	"github.com/Confialink/wallet-messages/internal/srvdiscovery"
	"context"
	"net/http"

	"github.com/Confialink/wallet-messages/internal/service/settings"

	notificationspb "github.com/Confialink/wallet-notifications/rpc/proto/notifications"
	"github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/inconshreveable/log15"
)

type NotificationService struct {
	logger log15.Logger
}

func NewNotificationService(logger log15.Logger) *NotificationService {
	return &NotificationService{logger: logger}
}

func (s *NotificationService) TriggerIncomingMessage(recipient *users.User, sender *users.User,
	messageId uint64, parentId *uint64) error {
	logger := s.logger.New("method", "TriggerNewMessage")
	client, err := s.getClient()
	if err != nil {
		logger.Error("failed to get pb client", "error", err)
		return err
	}

	var params *settings.MessageParams
	var entityId uint64
	if parentId != nil {
		params, err = s.getMessageParamsByMessageId(*parentId)
		entityId = *parentId
	} else {
		params, err = s.getMessageParamsByMessageId(messageId)
		entityId = messageId
	}

	if err != nil {
		logger.Error("failed to get message params", "error", err)
		return err
	}

	_, err = client.Dispatch(context.Background(), &notificationspb.Request{
		To:        recipient.UID,
		EventName: "IncomingMessage",
		TemplateData: &notificationspb.TemplateData{
			PrivateMessageAuthor:           sender.FirstName + " " + sender.LastName,
			PrivateMessageRecipient:        recipient.FirstName + " " + recipient.LastName,
			PrivateMessageUrl:              params.IncomingMessageUrl,
			PrivateMessageRecipientEditUrl: params.ProfileSettingsUrl,
			EntityType:                     "message",
			EntityID:                       entityId,
			MessageUnreadedCount:           1,
		},
	})

	return err
}

func (s *NotificationService) TriggerIncomingMessageAdmins(sender *users.User,
	messageId uint64, parentId *uint64) error {
	logger := s.logger.New("method", "TriggerNewMessageForAdmins")
	client, err := s.getClient()
	if err != nil {
		logger.Error("failed to get pb client", "error", err)
		return err
	}

	var params *settings.MessageParams
	var entityId uint64
	if parentId != nil {
		params, err = s.getMessageParamsByMessageId(*parentId)
		entityId = *parentId
	} else {
		params, err = s.getMessageParamsByMessageId(messageId)
		entityId = messageId
	}

	if err != nil {
		logger.Error("failed to get message params", "error", err)
		return err
	}

	_, err = client.Dispatch(context.Background(), &notificationspb.Request{
		EventName: "IncomingMessageAdmins",
		TemplateData: &notificationspb.TemplateData{
			PrivateMessageAuthor:           sender.FirstName + " " + sender.LastName,
			PrivateMessageUrl:              params.IncomingMessageUrl,
			PrivateMessageRecipientEditUrl: params.ProfileSettingsUrl,
			EntityType:                     "message",
			EntityID:                       entityId,
			MessageUnreadedCount:           1,
		},
	})

	return err
}

func (s *NotificationService) getClient() (notificationspb.NotificationHandler, error) {
	notificationsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameNotifications)
	if nil != err {
		return nil, err
	}
	return notificationspb.NewNotificationHandlerProtobufClient(notificationsUrl.String(), http.DefaultClient), nil
}

func (s *NotificationService) getMessageParamsByMessageId(id uint64) (*settings.MessageParams, error) {
	params, err := settings.GetMessageParamsByMessageId(id)
	if err != nil {
		return nil, err
	}

	return params, nil
}
