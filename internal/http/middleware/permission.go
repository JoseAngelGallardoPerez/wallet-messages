package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-messages/internal/auth"
	"github.com/Confialink/wallet-messages/internal/errcodes"
	"github.com/Confialink/wallet-messages/internal/http/handlers"
)

type PermissionChecker struct {
	authService *auth.Service
}

func NewPermissionChecker(authService *auth.Service) *PermissionChecker {
	return &PermissionChecker{authService}
}

func (s *PermissionChecker) Can(action, resource string) func(*gin.Context) {
	return func(c *gin.Context) {
		user := handlers.GetCurrentUser(c)
		if user == nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		if !s.authService.Can(user.RoleName, action, resource) {
			errcodes.AddError(c, errcodes.CodeForbidden)
			c.Abort()
			return
		}
	}
}

func (s *PermissionChecker) CanDynamic(action, resource string) func(*gin.Context) {
	return func(c *gin.Context) {
		user := handlers.GetCurrentUser(c)
		if user == nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		if !s.authService.CanDynamic(user, action, resource) {
			errcodes.AddError(c, errcodes.CodeForbidden)
			c.Abort()
			return
		}
	}
}

func (s *PermissionChecker) CanUpdateMessage() func(*gin.Context) {
	return func(c *gin.Context) {
		user := handlers.GetCurrentUser(c)
		if user == nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		message := handlers.GetRequestedMessage(c)
		if message == nil {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		if !s.authService.Can(user.RoleName, auth.ActionUpdate, auth.ResourceMessages) &&
			(user.UID != *message.SenderId ||
				user.RoleName != auth.RoleAdmin) {
			errcodes.AddError(c, errcodes.CodeForbidden)
			c.Abort()
			return
		}
	}
}

func (s *PermissionChecker) CanDeleteForMe() func(*gin.Context) {
	return func(c *gin.Context) {
		user := handlers.GetCurrentUser(c)
		if user == nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		message := handlers.GetRequestedMessage(c)
		if message == nil {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		if !s.authService.Can(user.RoleName, auth.ActionDelete, auth.ResourceMessages) ||
			(message.SenderId != nil && user.UID != *message.SenderId &&
				(message.RecipientId != nil && user.UID != *message.RecipientId)) {
			errcodes.AddError(c, errcodes.CodeForbidden)
			c.Abort()
			return
		}
	}
}

func (s *PermissionChecker) CanDeleteForAll() func(*gin.Context) {
	return func(c *gin.Context) {
		user := handlers.GetCurrentUser(c)
		if user == nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		message := handlers.GetRequestedMessage(c)
		if message == nil {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		if !s.authService.Can(user.RoleName, auth.ActionDelete, auth.ResourceMessages) ||
			message.SenderId != nil && user.UID != *message.SenderId {
			errcodes.AddError(c, errcodes.CodeForbidden)
			c.Abort()
			return
		}
	}
}
