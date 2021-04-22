package routes

import (
	"net/http"

	"github.com/Confialink/wallet-messages/internal/authentication"
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/Confialink/wallet-pkg-service_names"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-messages/internal/auth"
	"github.com/Confialink/wallet-messages/internal/config"
	"github.com/Confialink/wallet-messages/internal/dao"
	"github.com/Confialink/wallet-messages/internal/di"
	"github.com/Confialink/wallet-messages/internal/http/handlers"
	"github.com/Confialink/wallet-messages/internal/http/middleware"
	"github.com/Confialink/wallet-messages/internal/version"
)

var r *gin.Engine
var c = di.Container

func initRoutes() {
	r = gin.New()

	var cfg *config.Config
	var logger log15.Logger
	var daoMessage *dao.Message
	var authService *auth.Service
	err := c.Invoke(func(config *config.Config, l log15.Logger, d *dao.Message, authServ *auth.Service) {
		cfg = config
		logger = l
		daoMessage = d
		authService = authServ
	})
	if err != nil {
		panic(err)
	}

	// Middleware

	mwAuth := authentication.Middleware(logger.New("Middleware", "Authentication"))
	mwCors := middleware.CorsMiddleware(cfg.Cors)
	mwAdminOrRoot := middleware.AdminOrRoot()
	mwRequestedMessage := middleware.RequestedMessage(daoMessage)
	mwPermissons := middleware.NewPermissionChecker(authService)
	hndlrMessage := handlers.Factory.MessageHandlerFactory()

	// Routes

	r.GET("/messages/health-check", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/messages/build", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.BuildInfo)
	})

	apiGroup := r.Group(service_names.Messages.Internal, mwCors)
	apiGroup.Use(
		gin.Recovery(),
		gin.Logger(),
		errors.ErrorHandler(logger.New("Middleware", "Errors")),
	)

	privateGroup := apiGroup.Group("/private")
	{
		v1Group := privateGroup.Group("/v1", mwAuth)
		{
			v1Group.GET("/messages/:id", hndlrMessage.GetHandler)
			v1Group.GET("/messages", mwPermissons.Can(auth.ActionRead, auth.ResourceMessages), hndlrMessage.ListHandler)
			v1Group.POST("/messages", mwPermissons.CanDynamic(auth.ActionCreate, auth.ResourceMessages), hndlrMessage.CreateHandler)
			v1Group.DELETE("/messages/:id/for-me", mwRequestedMessage, mwPermissons.CanDeleteForMe(), hndlrMessage.DeleteForMeHandler)
			v1Group.DELETE("/messages/:id/for-all", mwRequestedMessage, mwPermissons.CanDeleteForAll(), hndlrMessage.DeleteForAllHandler)

			adminGroup := v1Group.Group("/admin", mwAdminOrRoot)
			{
				adminGroup.GET("/messages/unassigned-and-incoming", mwPermissons.Can(auth.ActionRead, auth.ResourceMessagesAdmin), hndlrMessage.ListUnassignedAndIncomingHandler)
				adminGroup.GET("/messages", mwPermissons.Can(auth.ActionRead, auth.ResourceMessagesAdmin), hndlrMessage.AdminListHandler)
				adminGroup.POST("/messages/send-to-all", mwPermissons.CanDynamic(auth.ActionCreate, auth.ResourceMessages), hndlrMessage.SendToAllHandler)
				adminGroup.POST("/messages/send-to-user-group", mwPermissons.CanDynamic(auth.ActionCreate, auth.ResourceMessages), hndlrMessage.SendToUserGroupHandler)
				adminGroup.POST("/messages/send-to-specific-users", mwPermissons.CanDynamic(auth.ActionCreate, auth.ResourceMessages), hndlrMessage.SendToSpecificUsersHandler)
			}

			csvGroup := v1Group.Group("/csv")
			{
				csvGroup.GET("/messages", mwPermissons.Can(auth.ActionRead, auth.ResourceMessages), hndlrMessage.CsvListHandler)

				adminGroup := csvGroup.Group("/admin", mwAdminOrRoot, mwPermissons.Can(auth.ActionRead, auth.ResourceMessagesAdmin))
				{
					adminGroup.GET("/messages/unassigned-and-incoming", hndlrMessage.CsvListUnassignedAndIncomingHandler)
					adminGroup.GET("/messages", hndlrMessage.CsvAdminListHandler)
				}
			}

			countGroup := v1Group.Group("/count")
			{
				messagesGroup := countGroup.Group("/messages")
				{
					messagesGroup.GET("/unread", mwPermissons.Can(auth.ActionRead, auth.ResourceMessages), hndlrMessage.CountUnreadMessagesHandler)
				}
			}
		}
	}

	// Handle OPTIONS request
	r.OPTIONS("/*cors", hndlrMessage.OptionsHandler, mwCors)

	r.NoRoute(hndlrMessage.NotFoundHandler)
}

func GetRouter() *gin.Engine {
	if nil == r {
		initRoutes()
	}

	return r
}
