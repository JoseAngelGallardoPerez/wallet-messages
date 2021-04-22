package handlers

import (
	"github.com/Confialink/wallet-messages/internal/auth"
	"github.com/Confialink/wallet-messages/internal/dao"
	"github.com/Confialink/wallet-messages/internal/di"
	"github.com/Confialink/wallet-messages/internal/service"
	"github.com/inconshreveable/log15"
)

var c = di.Container

// Factory struct for functory methods
var Factory *factory

type factory struct{}

func init() {
	Factory = &factory{}
}

func (f *factory) MessageHandlerFactory() *MessageHandler {
	var h *MessageHandler
	err := c.Invoke(func(
		d *dao.Message,
		a *auth.Service,
		u *service.UserService,
		n *service.NotificationService,
		c *service.Csv,
		m *service.Message,
		l log15.Logger,
	) {
		h = NewMessageHandler(d, a, u, n, c, m, l)
	})
	if err != nil {
		panic(err)
	}

	return h
}
