package handlers

import (
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/Confialink/wallet-pkg-utils/csv"
	"github.com/gin-gonic/gin"
)

// CsvListHandler returns the list of messages
func (self *MessageHandler) CsvListHandler(c *gin.Context) {
	currentUser := self.mustGetCurrentUser(c)

	c.Request.URL.Query().Set("limit", "10000")
	items, err := self.dao.FindByUserAndParams(currentUser.UID, c.Request.URL.Query())

	self.messageService.FillRecipientsAndSenders(items)

	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	file, err := self.csvService.MessagesToCsv(items, currentUser)
	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	if err := csv.Send(file, c.Writer); err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
	}
}

// CsvAdminListHandler returns the list of messages for admin
func (self *MessageHandler) CsvAdminListHandler(c *gin.Context) {
	currentUser := self.mustGetCurrentUser(c)

	c.Request.URL.Query().Set("limit", "10000")
	items, err := self.dao.FindByParams(c.Request.URL.Query(), currentUser.UID)

	self.messageService.FillRecipientsAndSenders(items)

	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	file, err := self.csvService.MessagesToCsv(items, currentUser)
	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	if err = csv.Send(file, c.Writer); err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
	}
}

// CsvListUnassignedAndIncomingHandler returns the list of unassigned messages
func (self *MessageHandler) CsvListUnassignedAndIncomingHandler(c *gin.Context) {
	c.Request.URL.Query().Set("limit", "10000")
	items, err := self.dao.FindUnassignedAndIncoming(c.Request.URL.Query(), self.mustGetCurrentUser(c).UID)

	self.messageService.FillRecipientsAndSenders(items)

	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	file, err := self.csvService.MessagesToCsv(items, self.mustGetCurrentUser(c))
	if err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
		return
	}

	if err := csv.Send(file, c.Writer); err != nil {
		errors.AddErrors(c, &errors.PrivateError{Message: err.Error()})
	}
}
