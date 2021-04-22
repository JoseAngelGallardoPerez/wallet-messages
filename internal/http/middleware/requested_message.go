package middleware

import (
	"net/http"

	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-messages/internal/dao"
	"github.com/Confialink/wallet-messages/internal/errcodes"
	"github.com/Confialink/wallet-messages/internal/http/handlers"
)

// CorsMiddleware cors middleware
func RequestedMessage(dao *dao.Message) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := handlers.GetCurrentUser(ctx)
		if user == nil {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		id, typedErr := handlers.GetIdParam(ctx)
		if typedErr != nil {
			errors.AddErrors(ctx, typedErr)
			ctx.Abort()
			return
		}

		message, err := dao.FindByID(id, user.UID, true)
		if err != nil {
			errcodes.AddError(ctx, errcodes.CodeMessageNotFound)
			ctx.Abort()
			return
		}

		ctx.Set("_requested_message", message)
	}
}
