package main

import (
	"fmt"
	"net/http"
	"os"

	gateway "spec-commentor/ms/gateway"

	// "spec-commentor/pkg/auth"
	// "spec-commentor/pkg/clock"
	// "spec-commentor/pkg/grpc"
	"spec-commentor/pkg/logger"

	"spec-commentor/config"
	"github.com/rs/zerolog/log"

	"spec-commentor/ms/gateway/generated"
)

func main() {
	if err := run(); err != nil {
		log.Error().Err(err).Msg("shutting down")
		os.Exit(1)
	}
}

func run() error {
	// Configuration
	cfg := config.NewGWServiceConfig()
	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	// Logging
	logger, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		return err
	}

	// jwtAuthenticator, err := auth.NewJwtAuthenticator(cfg.JWT, clock.New())
	// if err != nil {
	// 	logger.Error().Msg("Error initialization of the JWT Authenticator")
	// 	return err
	// }

	// // Authentication service grpc client
	// cfgAuthentication := config.NewAuthenticationServiceConfig()
	// authenticationConn := grpc.NewGRPCConnect(cfgAuthentication.GRPC, logger, jwtAuthenticator)
	// defer authenticationConn.Close()
	// authenticationClient := pbAuthentication.NewAuthenticationAPIServiceClient(authenticationConn)

	// // MarketPlace service grpc client
	// cfgMarketPlace := config.NewMarketPlaceServiceConfig()
	// marketplaceConn := grpc.NewGRPCConnect(cfgMarketPlace.GRPC, logger, jwtAuthenticator)
	// defer marketplaceConn.Close()
	// marketplaceClient := pbMarketPlace.NewMarketPlaceAPIServiceClient(marketplaceConn)

	server := gateway.NewServer(logger, cfg)
	options := gateway.NewServerOptions(&server)
	router := generated.HandlerWithOptions(server, options)

	logger.Info().Msgf("Start Admin GateWay server: %s", addr)

	err = http.ListenAndServeTLS(addr, "localhost.crt", "localhost.key", router)
	if err != nil {
		logger.Error().Msg("Error serving http")
	}

	return err
}
