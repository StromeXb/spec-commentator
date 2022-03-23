package grpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func StartGRPCServer(cfg *Config, server *grpc.Server, logger *zerolog.Logger) error {
	// Make a channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start gRPC server listening for requests.
	grpcListener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logger.Error().Err(err).Str("port", cfg.Port).Str("host", cfg.Host).Msg("Failed listen on gRPC port")

		return err
	}
	go func() {
		serverErrors <- server.Serve(grpcListener)
	}()
	logger.Info().Str("port", cfg.Port).Str("host", cfg.Host).Msg("GRPC service started")

	// Blocking and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("running server: %w", err)
	case sig := <-shutdown:
		logger.Info().Str("signal", sig.String()).Msg("Start shutdown")
		server.GracefulStop()
	}

	return nil
}

func NewServer(log *zerolog.Logger) *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			RecoveryUnaryInterceptor(log),
		),
	)
}
