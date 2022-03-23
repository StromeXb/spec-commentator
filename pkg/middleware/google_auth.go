package middleware

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"spec-commentor/pkg/auth"

	"github.com/rs/zerolog"
)

const (
	Userinfo key = 0
)

type key int

// AuthMiddleware creates middleware that checks authorization for the endpoint according to the
// logics that implemented in the Authenticator object.
// In case of successful authorization, claims will be added to the request's context under the UserCredentials key.
// Also AuthMiddleware will add the JWT token under the JWTCredentials key, that should be used for cross-service authorization.
func GoogleAuthMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token := r.Header.Get("Access-Token")
			if token == "" {
				http.Error(w, "token not provaided", http.StatusUnauthorized)
				return
			}
			logger.Info().Msg(token)
			url := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token
			resp, err := http.Get(url)
			if err != nil {
				logger.Debug().Err(err).Msgf("Could not get response %s\n", err.Error())
				http.Error(w, fmt.Sprintf("error checking auth: %s", err.Error()), http.StatusUnauthorized)
				return
			}
			defer resp.Body.Close()
			content, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Debug().Err(err).Msgf("Could not read response %s\n", err.Error())
				http.Error(w, fmt.Sprintf("error checking auth: %s", err.Error()), http.StatusUnauthorized)
			}

			ctx = context.WithValue(ctx, Userinfo, string(content))
			r = r.WithContext(ctx)

			// serve
			next.ServeHTTP(w, r)
		})
	}
}

func GetGoogleUserInfoFromContext(ctx context.Context) *int64 {
	cred := ctx.Value(Userinfo)
	if cred == nil {
		return nil
	}
	return cred.(*auth.UserCredentials).Id
}
