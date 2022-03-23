package grpc

import (
	"context"
	"fmt"
	"runtime/debug"

	"spec-commentor/pkg/auth"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

const (
	maxCallRecvMsgSize = 1024 * 1024 * 20
	Host               = "127.0.0.1"
)

type Config struct {
	Host       string
	Port       string
	Reflection bool
}

func RecoveryUnaryInterceptor(logger *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer func() {
			if info := recover(); info != nil {
				logger.Error().
					Interface("recover_info", info).
					Bytes("debug_stack", debug.Stack()).
					Msg("panic_on_request")
			}
		}()

		return handler(ctx, req)
	}
}

func NewGRPCConnect(cfg *Config, logger *zerolog.Logger, jwtAuth auth.JwtAuthenticator) *grpc.ClientConn {
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	logger.Info().Msgf("GRPC client address connect: %s", address)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize)))
	opts = append(opts,
		grpc.WithChainUnaryInterceptor(
			JWTAuthorizationClientUnaryInterceptor(jwtAuth, logger),
		),
	)
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logger.Error().Err(err).Msgf("Fail GRPC client address connect: %s", address)
	}

	return conn
}
