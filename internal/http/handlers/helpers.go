package handlers

import (
	"strconv"

	"github.com/Confialink/wallet-pkg-errors"
	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-messages/internal/errcodes"
	"github.com/Confialink/wallet-messages/internal/model"
)

// getIdParam returns id or nil
func GetIdParam(c *gin.Context) (uint64, errors.TypedError) {
	id := c.Params.ByName("id")

	// convert string to uint
	id64, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		return 0, errcodes.CreatePublicError(errcodes.CodeInvalidQueryParameters, "id param must be an integer")
	}

	return uint64(id64), nil
}

// getCurrentUser returns current user or nil
func GetCurrentUser(c *gin.Context) *userpb.User {
	user, ok := c.Get("_user")
	if !ok {
		return nil
	}
	return user.(*userpb.User)
}

// mw already get message by UI from DB. retrieve requested message from gin context
func GetRequestedMessage(ctx *gin.Context) *model.Message {
	message, ok := ctx.Get("_requested_message")
	if !ok {
		return nil
	}
	return message.(*model.Message)
}
