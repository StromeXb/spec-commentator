package grpc

import (
	"context"
	"fmt"

	"spec-commentor/pkg/auth"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	AuthMdKey       = "authorization"
	TokenInfoMethod = "/authentication.proto.api.AuthenticationAPIService/TokenInfo"
)

// JWTAuthorizationServerInterceptor unary interceptor function to handle authorize per RPC call
func JWTAuthorizationServerUnaryInterceptor(
	authenticator auth.JwtAuthenticator,
	logger *zerolog.Logger,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == TokenInfoMethod {
			return handler(ctx, req)
		}
		// получаем токен из контекста
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Error()
			return nil, fmt.Errorf("retrieving metadata is failed")
		}

		authHeader, ok := md[AuthMdKey]
		if !ok {
			return nil, fmt.Errorf("authorization token is not supplied")
		}
		token := authHeader[0]

		// раскодируем и валидируем токен
		cred, err := authenticator.ParseToken(token)
		if err != nil {
			return nil, err
		}

		// добавляем данные о пользователе в контекст
		ctx = context.WithValue(ctx, auth.JwtAuthContextKey, cred)
		ctx = context.WithValue(ctx, auth.JwtAuthCredKey, token)

		return handler(ctx, req)
	}
}

// JWTAuthorizationClientUnaryInterceptor gets JWT from JWTCredentials in the context and adds it to the metadata of a grpc-request
func JWTAuthorizationClientUnaryInterceptor(
	authenticator auth.JwtAuthenticator,
	logger *zerolog.Logger,
) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if method == TokenInfoMethod {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		// получаем JWT из контекста текущего запроса
		value := ctx.Value(auth.JwtAuthCredKey)
		if value == nil {
			return fmt.Errorf("request has empty user JWT credentials")
		}

		// добавляем JWT в метаданные запроса
		ctx = metadata.AppendToOutgoingContext(ctx, AuthMdKey, value.(string))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
