package service

import (
	"github.com/Confialink/wallet-messages/internal/errcodes"
	"github.com/Confialink/wallet-messages/internal/model"
	userModel "github.com/Confialink/wallet-users/rpc/proto/users"
)

type Message struct {
	userService *UserService
}

func NewMessage(userService *UserService) *Message {
	return &Message{
		userService: userService,
	}
}

func (self *Message) FillRecipientsAndSenders(messages []*model.Message) error {
	err := self.FillRecipients(messages)
	if nil != err {
		return err
	}

	err = self.FillSenders(messages)
	if nil != err {
		return err
	}

	return nil
}

func (self *Message) FillRecipientAndSender(message *model.Message) error {
	if nil != message.RecipientId {
		recipient, err := self.userService.GetByUID(*message.RecipientId)
		if nil != err {
			return err
		}

		self.fillRecipient(message, recipient)
	}

	if nil != message.SenderId {
		sender, err := self.userService.GetByUID(*message.SenderId)
		if nil != err {
			return err
		}

		self.fillSender(message, sender)
	}
	return nil
}

func (self *Message) FillRecipient(message *model.Message) error {
	if nil != message.RecipientId {
		recipient, err := self.userService.GetByUID(*message.RecipientId)
		if nil != err {
			return errcodes.CreatePublicError(errcodes.CodeUserNotFound)
		}

		self.fillRecipient(message, recipient)
	}

	return nil
}

func (self *Message) FillSender(message *model.Message) error {
	if nil != message.SenderId {
		sender, err := self.userService.GetByUID(*message.SenderId)
		if nil != err {
			return errcodes.CreatePublicError(errcodes.CodeUserNotFound)
		}

		self.fillSender(message, sender)
	}
	return nil
}

func (self *Message) FillRecipients(messages []*model.Message) error {
	recipientsIds := make([]string, 0)
	for _, v := range messages {
		if nil != v.RecipientId && !self.isExist(recipientsIds, *v.RecipientId) {
			recipientsIds = append(recipientsIds, *v.RecipientId)
		}
	}

	recipients, err := self.userService.GetByUIDs(recipientsIds)
	if err != nil {
		return errcodes.CreatePublicError(errcodes.CodeUserNotFound)
	}

	for _, v := range messages {
		if nil != v.RecipientId {
			recipient := self.findUserById(recipients, *v.RecipientId)
			if recipient != nil {
				self.fillRecipient(v, recipient)
			}
		}
	}

	return nil
}

func (self *Message) FillSenders(messages []*model.Message) error {
	sendersIds := make([]string, 0)
	for _, v := range messages {
		if nil != v.SenderId && !self.isExist(sendersIds, *v.SenderId) {
			sendersIds = append(sendersIds, *v.SenderId)
		}
	}

	senders, err := self.userService.GetByUIDs(sendersIds)
	if err != nil {
		return errcodes.CreatePublicError(errcodes.CodeUserNotFound)
	}

	for _, v := range messages {
		if nil != v.SenderId {
			sender := self.findUserById(senders, *v.SenderId)
			if sender != nil {
				self.fillSender(v, sender)
			}
		}
	}

	return nil
}

func (self *Message) fillRecipient(
	message *model.Message, recipient *userModel.User,
) {
	message.Recipient = &model.User{
		UID:       &recipient.UID,
		Email:     &recipient.Email,
		Username:  &recipient.Username,
		FirstName: &recipient.FirstName,
		LastName:  &recipient.LastName,
		RoleName:  &recipient.RoleName,
		GroupId:   &recipient.GroupId,
	}
}

func (self *Message) fillSender(
	message *model.Message, sender *userModel.User,
) {
	message.Sender = &model.User{
		UID:       &sender.UID,
		Email:     &sender.Email,
		Username:  &sender.Username,
		FirstName: &sender.FirstName,
		LastName:  &sender.LastName,
		RoleName:  &sender.RoleName,
		GroupId:   &sender.GroupId,
	}
}

func (self *Message) findUserById(
	array []*userModel.User, id string,
) *userModel.User {
	for _, v := range array {
		if v.UID == id {
			return v
		}
	}
	return nil
}

func (self *Message) isExist(array []string, elem string) bool {
	for _, v := range array {
		if v == elem {
			return true
		}
	}
	return false
}
