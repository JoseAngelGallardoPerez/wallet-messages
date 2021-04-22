package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/Confialink/wallet-pkg-utils/pointer"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-messages/internal/auth"
	"github.com/Confialink/wallet-messages/internal/dao"
	"github.com/Confialink/wallet-messages/internal/errcodes"
	"github.com/Confialink/wallet-messages/internal/http/response"
	"github.com/Confialink/wallet-messages/internal/model"
	"github.com/Confialink/wallet-messages/internal/service"
)

// Service
type MessageHandler struct {
	dao                 *dao.Message
	authService         *auth.Service
	userService         *service.UserService
	notificationService *service.NotificationService
	csvService          *service.Csv
	messageService      *service.Message
	logger              log15.Logger
}

// NewMessageHandler creates new message handler
func NewMessageHandler(
	dao *dao.Message,
	authService *auth.Service,
	userService *service.UserService,
	notificationService *service.NotificationService,
	csvService *service.Csv,
	messageService *service.Message,
	logger log15.Logger,
) *MessageHandler {
	return &MessageHandler{dao, authService, userService, notificationService, csvService, messageService, logger}
}

// OptionsHandler handle options request
func (self *MessageHandler) OptionsHandler(c *gin.Context) {}

// ListHandler returns the list of messages
func (self *MessageHandler) ListHandler(c *gin.Context) {
	currentUser := self.mustGetCurrentUser(c)

	items, err := self.dao.FindByUserAndParams(currentUser.UID, c.Request.URL.Query())
	count, err := self.dao.CountByUserAndParams(currentUser.UID, c.Request.URL.Query())

	self.messageService.FillRecipientsAndSenders(items)

	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}
	r, _ := response.NewResponseWithListAndLinks(items, c, *count)
	c.JSON(http.StatusOK, r)
}

// AdminListHandler returns the list of messages for admin
func (self *MessageHandler) AdminListHandler(c *gin.Context) {
	currentUser := self.mustGetCurrentUser(c)

	items, err := self.dao.FindByParams(c.Request.URL.Query(), currentUser.UID)
	count, err := self.dao.CountByParams(c.Request.URL.Query(), currentUser.UID)

	self.messageService.FillRecipientsAndSenders(items)

	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}
	r, _ := response.NewResponseWithListAndLinks(items, c, *count)
	c.JSON(http.StatusOK, r)
}

// ListHandler returns the list of unassigned messages
func (self *MessageHandler) ListUnassignedAndIncomingHandler(c *gin.Context) {
	currentUser := self.mustGetCurrentUser(c)
	items, err := self.dao.FindUnassignedAndIncoming(c.Request.URL.Query(), currentUser.UID)
	count, err := self.dao.CountUnassignedAndIncoming(c.Request.URL.Query(), currentUser.UID)

	self.messageService.FillRecipientsAndSenders(items)

	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}
	r, _ := response.NewResponseWithListAndLinks(items, c, *count)
	c.JSON(http.StatusOK, r)
}

// GetHandler returns message by id
func (self *MessageHandler) GetHandler(c *gin.Context) {
	id, typedErr := GetIdParam(c)
	if typedErr != nil {
		errors.AddErrors(c, typedErr)
		return
	}

	currentUser := self.mustGetCurrentUser(c)
	var message *model.Message

	var err error
	if currentUser.RoleName == "root" || currentUser.RoleName == "admin" {
		message, err = self.dao.FindByIDWithChildrenForAdmin(id, currentUser.UID)
	} else {
		message, err = self.dao.FindByIDWithChildren(id, currentUser.UID)
	}

	if err != nil {
		self.logger.Error(err.Error())
		errcodes.AddError(c, errcodes.CodeMessageNotFound)
		return
	}

	if !self.authService.Can(
		self.mustGetCurrentUser(c).RoleName,
		auth.ActionRead,
		auth.ResourceMessages,
	) &&
		(*message.SenderId != currentUser.UID ||
			message.RecipientId != &currentUser.UID) {
		errcodes.AddError(c, errcodes.CodeForbidden)
		return
	}

	if (currentUser.RoleName == "root" || currentUser.RoleName == "admin") && message.RecipientId == nil {
		message.IsRecipientRead = true
	}

	// mark as read if sender or recipient is a current user
	// or show and delete if message marked with flag DeleteAfterRead
	if nil != message.SenderId && *message.SenderId == currentUser.UID {
		message.IsSenderRead = true
	} else if nil != message.RecipientId && *message.RecipientId == currentUser.UID &&
		(nil == message.DeleteAfterRead ||
			*message.DeleteAfterRead != true) {
		message.IsRecipientRead = true
	} else if nil != message.RecipientId && *message.RecipientId == currentUser.UID &&
		(nil != message.DeleteAfterRead &&
			*message.DeleteAfterRead == true) {
		self.showAndDeleteMessage(*c, *message)
		return
	}

	if nil != message.RecipientId || (currentUser.RoleName == "root" || currentUser.RoleName == "admin") && message.RecipientId == nil {
		message, err = self.dao.Update(message)
	}

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	self.messageService.FillRecipientAndSender(message)

	for _, chld := range message.Children {
		self.messageService.FillRecipientAndSender(chld)
	}

	c.JSON(http.StatusOK, response.NewResponse().SetData(message))
}

// CreateHandler creates new message
func (self *MessageHandler) CreateHandler(c *gin.Context) {
	currentUser := self.mustGetCurrentUser(c)

	var public model.MessagePublic

	if err := c.ShouldBindJSON(&public); err != nil {
		errors.AddShouldBindError(c, err)
		return
	}

	publicJson, err := json.Marshal(public)

	if nil != err {
		errors.AddShouldBindError(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	var message *model.Message
	json.Unmarshal(publicJson, &message)
	message.SenderId = &self.mustGetCurrentUser(c).UID

	if nil == message.ParentId {
		message.IsRecipientIncoming = true
		message.IsSenderRead = true
	}

	if nil != message.ParentId {
		parent, err := self.dao.FindByID(*message.ParentId, currentUser.UID, false)

		if nil != err {
			errcodes.AddError(c, errcodes.CodeMessageNotFound)
			return
		}

		if (currentUser.RoleName == "root" || currentUser.RoleName == "admin") && parent.RecipientId == nil {
			parent.RecipientId = &currentUser.UID
		}

		if nil != parent.RecipientId && *parent.SenderId == self.mustGetCurrentUser(c).UID {
			parent.IsSenderRead = true
			parent.IsRecipientRead = false
			parent.IsRecipientIncoming = true
		} else if nil != parent.RecipientId && *parent.RecipientId == self.mustGetCurrentUser(c).UID {
			parent.IsRecipientRead = true
			parent.IsSenderRead = false
			parent.IsRecipientIncoming = false
		} else {
			parent.IsRecipientRead = false
		}

		self.dao.Update(parent)
	}

	createdMessage, err := self.dao.Create(message)
	self.messageService.FillSender(createdMessage)

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	go self.triggerIncomingMessageNotification(createdMessage)

	c.Status(http.StatusOK)
	c.JSON(http.StatusCreated, response.NewResponse().SetData(createdMessage))
}

// SendToAllHandler send messages to all users
func (self *MessageHandler) SendToSpecificUsersHandler(c *gin.Context) {
	data := struct {
		Message         string   `json:"message" binding:"required"`
		Subject         *string  `json:"subject"`
		DeleteAfterRead *bool    `json:"deleteAfterRead"`
		UsersIds        []string `json:"usersIds"`
	}{}

	if err := c.ShouldBindJSON(&data); err != nil {
		errors.AddShouldBindError(c, err)
		return
	}

	var privateErrors []error

	_, err := self.userService.GetByUIDs(data.UsersIds)
	if err != nil {
		errcodes.AddError(c, errcodes.CodeUserNotFound)
		return
	}

	messages := make([]*model.Message, 0)
	for _, userId := range data.UsersIds {
		if userId != self.mustGetCurrentUser(c).UID {
			message := &model.Message{
				MessagePublic: model.MessagePublic{
					Message:         data.Message,
					Subject:         data.Subject,
					DeleteAfterRead: data.DeleteAfterRead,
					RecipientId:     &userId,
				},
				MessagePrivate: model.MessagePrivate{
					SenderId:            &self.mustGetCurrentUser(c).UID,
					IsSenderRead:        true,
					IsRecipientIncoming: true,
				},
			}

			message, err := self.dao.Create(message)
			messages = append(messages, message)

			if nil != err {
				privateErrors = append(privateErrors, err)
			} else {
				go self.triggerIncomingMessageNotification(message)
			}
		}
	}

	if len(privateErrors) > 0 {
		for _, err := range privateErrors {
			errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
	c.JSON(http.StatusCreated, response.NewResponse().SetData(messages))
}

// SendToAllHandler send messages to all users
func (self *MessageHandler) SendToAllHandler(c *gin.Context) {
	data := struct {
		Message         string  `json:"message" binding:"required"`
		Subject         *string `json:"subject"`
		DeleteAfterRead *bool   `json:"deleteAfterRead"`
	}{}

	if err := c.ShouldBindJSON(&data); err != nil {
		errors.AddShouldBindError(c, err)
		return
	}

	users, err := self.userService.GetAll()

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	messages := make([]*model.Message, 0)

	for _, user := range users {
		if user.UID != self.mustGetCurrentUser(c).UID {
			message := model.Message{
				MessagePublic: model.MessagePublic{
					Message:         data.Message,
					Subject:         data.Subject,
					DeleteAfterRead: data.DeleteAfterRead,
					RecipientId:     &user.UID,
				},
				MessagePrivate: model.MessagePrivate{
					SenderId:            &self.mustGetCurrentUser(c).UID,
					IsSenderRead:        true,
					IsRecipientIncoming: true,
				},
			}
			messages = append(messages, &message)
		}
	}

	messages, err = self.dao.BulkCreate(messages)

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	go self.triggerIncomingMessagesNotifications(messages)

	c.Status(http.StatusOK)
	c.JSON(http.StatusCreated, response.NewResponse().SetData(messages))
}

// SendToAllHandler send messages to all users
func (self *MessageHandler) SendToUserGroupHandler(c *gin.Context) {
	data := struct {
		Message         string  `json:"message" binding:"required"`
		Subject         *string `json:"subject"`
		DeleteAfterRead *bool   `json:"deleteAfterRead"`
		UserGroupId     uint64  `json:"userGroupId" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&data); err != nil {
		errors.AddShouldBindError(c, err)
		c.Abort()
		return
	}

	users, err := self.userService.GetByGroup(data.UserGroupId)

	if nil != err {
		errcodes.AddError(c, errcodes.CodeUserNotFound)
		return
	}

	messages := make([]*model.Message, 0)

	for _, user := range users {
		if user.UID != self.mustGetCurrentUser(c).UID {
			message := model.Message{
				MessagePublic: model.MessagePublic{
					Message:         data.Message,
					Subject:         data.Subject,
					DeleteAfterRead: data.DeleteAfterRead,
					RecipientId:     &user.UID,
				},
				MessagePrivate: model.MessagePrivate{
					SenderId:            &self.mustGetCurrentUser(c).UID,
					IsSenderRead:        true,
					IsRecipientIncoming: true,
				},
			}
			messages = append(messages, &message)
		}
	}

	messages, err = self.dao.BulkCreate(messages)

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	go self.triggerIncomingMessagesNotifications(messages)

	c.Status(http.StatusOK)
	c.JSON(http.StatusCreated, response.NewResponse().SetData(messages))
}

// DeleteForSenderHandler delete a message for sender
func (self *MessageHandler) DeleteForMeHandler(c *gin.Context) {
	message := GetRequestedMessage(c)
	if message == nil {
		errcodes.AddError(c, errcodes.CodeMessageNotFound)
		return
	}
	currentUser := self.mustGetCurrentUser(c)

	if message.RecipientId != nil && *message.RecipientId == currentUser.UID {
		message.DeletedForRecipient = pointer.ToBool(true)
	} else if message.SenderId != nil && *message.SenderId == currentUser.UID {
		message.DeletedForSender = pointer.ToBool(true)
	}
	message, err := self.dao.Update(message)

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewResponse())
}

// DeleteForSenderHandler delete a message for recipient and sender
func (self *MessageHandler) DeleteForAllHandler(c *gin.Context) {
	message := GetRequestedMessage(c)
	if message == nil {
		errcodes.AddError(c, errcodes.CodeMessageNotFound)
		return
	}

	message.DeletedForRecipient = pointer.ToBool(true)
	message.DeletedForSender = pointer.ToBool(true)
	message, err := self.dao.Update(message)

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewResponse())
}

// CountUnreadMessagesHandler returns count of unread messages for current user
func (self *MessageHandler) CountUnreadMessagesHandler(c *gin.Context) {
	currentUser := self.mustGetCurrentUser(c)

	var count *int64
	var err error

	if currentUser.RoleName == auth.RoleAdmin || currentUser.RoleName == auth.RoleRoot {
		count, err = self.dao.CountUnreadByAdmin(currentUser.UID)
	} else {
		count, err = self.dao.CountUnreadByUser(currentUser.UID)
	}

	if nil != err {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewResponse().SetData(count))
}

// NotFoundHandler returns 404 NotFound
func (self *MessageHandler) NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
}

func (self *MessageHandler) showAndDeleteMessage(c gin.Context, message model.Message) {
	messageBeforeDelete := message
	err := self.dao.Delete(&message)

	if nil != err {
		errors.AddErrors(&c, &errors.PrivateError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewResponse().SetData(messageBeforeDelete))
	return
}

func (self *MessageHandler) triggerIncomingMessagesNotifications(messages []*model.Message) {
	for _, v := range messages {
		self.triggerIncomingMessageNotification(v)
	}
}

func (self *MessageHandler) triggerIncomingMessageNotification(message *model.Message) {
	logger := self.logger.New("action", "triggerIncomingMessageNotification")

	if message.SenderId != nil {
		sender, err := self.userService.GetByUID(*message.SenderId)
		if err != nil {
			logger.Error("can't retrieve sender from user service", "err", err)
			return
		}

		if message.RecipientId != nil {
			recipient, err := self.userService.GetByUID(*message.RecipientId)
			if err != nil {
				logger.Error("can't retrieve recipient from user service", "err", err)
				return
			}

			err = self.notificationService.TriggerIncomingMessage(recipient, sender, message.ID, message.ParentId)
			if err != nil {
				logger.Error("can't send incoming message notification", "err", err)
			}
		} else {
			err = self.notificationService.TriggerIncomingMessageAdmins(sender, message.ID, message.ParentId)
			if err != nil {
				logger.Error("can't send incoming message notification", "err", err)
			}
		}
	}
}

// mustGetCurrentUser returns current user or throw error
func (self *MessageHandler) mustGetCurrentUser(c *gin.Context) *userpb.User {
	user := GetCurrentUser(c)
	if nil == user {
		panic("user must be set")
	}
	return user
}
