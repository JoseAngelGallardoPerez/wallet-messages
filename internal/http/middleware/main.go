package middleware

import (
	"net/http"

	"github.com/Confialink/wallet-messages/internal/errcodes"
	"github.com/Confialink/wallet-pkg-env_config"
	"github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware cors middleware
func CorsMiddleware(config *env_config.Cors) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowMethods = config.Methods
	for _, origin := range config.Origins {
		if origin == "*" {
			corsConfig.AllowAllOrigins = true
		}
	}
	if !corsConfig.AllowAllOrigins {
		corsConfig.AllowOrigins = config.Origins
	}
	corsConfig.AllowHeaders = config.Headers

	return cors.New(corsConfig)
}

func AdminOrRoot() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := ctx.Get("_user")
		if !ok {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		roleName := user.(*users.User).RoleName
		if roleName != "admin" && roleName != "root" {
			ctx.Status(http.StatusForbidden)
			_ = ctx.Error(errcodes.CreatePublicError(errcodes.CodeForbidden, "you are not allowed to perform this action"))
			ctx.Abort()
			return
		}
	}
}
