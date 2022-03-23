package config

import (
	"spec-commentor/pkg/auth"
	"spec-commentor/pkg/http"
	"spec-commentor/pkg/logger"
	"spec-commentor/pkg/utils"
	"github.com/spf13/viper"
)

const (
	envGwServerEnvironment = "GW_SERVER_ENVIRONMENT"
	envGwServerScheme      = "GW_SERVER_SCHEME"
	envGwServerHost        = "GW_SERVER_HOST"
	envGwServerPort        = "GW_SERVER_PORT"
)

type GWConfig struct {
	Logger            *logger.Config
	Server            *http.Config
	JWT               *auth.JWTConfig
	MarketplaceServer *http.Config
	ServerEnvironment utils.ServerEnvironment
}

func NewGWServiceConfig() *GWConfig {
	v := viper.New()

	v.SetDefault(envGwServerEnvironment, utils.DEV)
	v.SetDefault(envGwServerScheme, http.Scheme)
	v.SetDefault(envGwServerHost, http.Host)
	v.SetDefault(envGwServerPort, http.Port)

	return &GWConfig{
		Logger:            logger.NewLoggerConfig(v),
		Server:            NewMobileGWServiceConfig(v),
		ServerEnvironment: utils.ServerEnvironment(v.GetString(envGwServerEnvironment)),
	}
}


func NewMobileGWServiceConfig(v *viper.Viper) *http.Config {
	return &http.Config{
		Scheme: v.GetString(envGwServerScheme),
		Host:   v.GetString(envGwServerHost),
		Port:   v.GetString(envGwServerPort),
	}
}
