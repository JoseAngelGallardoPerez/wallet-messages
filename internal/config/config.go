package config

import (
	"os"

	"github.com/Confialink/wallet-pkg-env_config"
	"github.com/Confialink/wallet-pkg-env_mods"
	"github.com/inconshreveable/log15"
)

type Config struct {
	Env          string
	Db           *env_config.Db
	Port         string
	Cors         *env_config.Cors
	ProtoBufPort string
}

func FromEnv(logger log15.Logger) *Config {
	c := &Config{}
	readConfig(c, logger)
	return c
}

// readConfig reads configs from ENV variables
func readConfig(cfg *Config, logger log15.Logger) {
	cfg.Port = os.Getenv("VELMIE_WALLET_MESSAGES_PORT")
	cfg.Env = env_config.Env("ENV", env_mods.Development)
	cfg.ProtoBufPort = os.Getenv("VELMIE_WALLET_MESSAGES_PROTO_BUF_PORT")

	defaultConfigReader := env_config.NewReader("messages")
	cfg.Cors = defaultConfigReader.ReadCorsConfig()
	cfg.Db = defaultConfigReader.ReadDbConfig()
	validateConfig(cfg, logger)
}

func validateConfig(cfg *Config, logger log15.Logger) {
	validator := env_config.NewValidator(logger)
	validator.ValidateCors(cfg.Cors, logger)
	validator.ValidateDb(cfg.Db, logger)
	validator.CriticalIfEmpty(cfg.Port, "VELMIE_WALLET_MESSAGES_PORT", logger)
	validator.CriticalIfEmpty(cfg.ProtoBufPort, "VELMIE_WALLET_MESSAGES_PROTO_BUF_PORT", logger)
}
