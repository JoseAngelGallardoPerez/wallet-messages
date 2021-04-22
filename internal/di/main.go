package di

import (
	"log"

	"github.com/Confialink/wallet-messages/internal/auth"
	"github.com/Confialink/wallet-messages/internal/config"
	"github.com/Confialink/wallet-messages/internal/dao"
	"github.com/Confialink/wallet-messages/internal/db"
	"github.com/Confialink/wallet-messages/internal/service"
	"github.com/Confialink/wallet-messages/rpc"
	"github.com/inconshreveable/log15"
	"go.uber.org/dig"
)

var Container *dig.Container

func init() {
	Container = dig.New()

	providers := []interface{}{
		// *gorm.DB
		db.NewConnection,
		// log15.Logger
		log15.New,
	}

	providers = append(providers, config.Providers()...)
	providers = append(providers, dao.Providers()...)
	providers = append(providers, service.Providers()...)
	providers = append(providers, auth.Providers()...)
	providers = append(providers, rpc.Providers()...)

	for _, provider := range providers {
		err := Container.Provide(provider)
		if err != nil {
			log.Fatal(err)
			panic("unable to init container: " + err.Error())
		}
	}
}
