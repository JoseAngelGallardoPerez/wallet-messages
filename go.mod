module github.com/Confialink/wallet-messages

go 1.13

replace github.com/Confialink/wallet-messages/rpc/messages => ./rpc/messages

require (
	github.com/Confialink/wallet-messages/rpc/messages v0.0.0-00010101000000-000000000000
	github.com/Confialink/wallet-notifications/rpc/proto/notifications v0.0.0-20210218064438-818cea3b20db
	github.com/Confialink/wallet-permissions/rpc/permissions v0.0.0-20210218072732-21caf4a66e86
	github.com/Confialink/wallet-pkg-discovery/v2 v2.0.0-20210217105157-30e31661c1d1
	github.com/Confialink/wallet-pkg-env_config v0.0.0-20210217112253-9483d21626ce
	github.com/Confialink/wallet-pkg-env_mods v0.0.0-20210217112432-4bda6de1ee2c
	github.com/Confialink/wallet-pkg-errors v0.1.1
	github.com/Confialink/wallet-pkg-service_names v0.0.0-20210217112604-179d69540dea
	github.com/Confialink/wallet-pkg-utils v0.0.0-20210217112822-e79f6d74cdc1
	github.com/Confialink/wallet-settings/rpc/proto/settings v0.0.0-20210218070334-b4153fc126a0
	github.com/Confialink/wallet-users/rpc/proto/users v0.0.0-20210218071418-0600c0533fb2
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/jinzhu/gorm v1.9.15
	github.com/kildevaeld/go-acl v0.0.0-20171228130000-7799b11f4759
	github.com/stretchr/testify v1.7.0
	go.uber.org/dig v1.10.0
)
