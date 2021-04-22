package main

import (
	"log"

	"github.com/Confialink/wallet-messages/internal/config"
	"github.com/Confialink/wallet-messages/internal/di"
	"github.com/Confialink/wallet-messages/internal/http/routes"
	"github.com/Confialink/wallet-messages/rpc"
	"github.com/Confialink/wallet-pkg-env_mods"
	"github.com/gin-gonic/gin"
)

// main: main function
func main() {
	c := di.Container

	var cfg *config.Config
	var pbServer *rpc.PbServer
	err := c.Invoke(func(c *config.Config, pb *rpc.PbServer) {
		cfg = c
		pbServer = pb
		gin.SetMode(env_mods.GetMode(c.Env))
	})
	if err != nil {
		panic(err)
	}

	ginRouter := routes.GetRouter()

	log.Printf("Starting API on port: %s", cfg.Port)

	// Start proto buf server
	go pbServer.Start()

	// Start gin server
	err = ginRouter.Run(":" + cfg.Port)
	if err != nil {
		panic(err)
	}
}
