package auth

import (
	"context"
	"errors"
	"net/http"
)

const (
	UserRoleAnonymous = "anonymous"
	UserRoleUser      = "user"
	UserRoleAdmin     = "admin"
)

var ErrAuthUnauthorized = errors.New("unauthorized")

const (
	JwtAuthCredKey    = AuthJWTCredentials("JWTCredentials")
	JwtAuthContextKey = AuthContext("UserCredentials")
)

type (
	UserRole           string
	AuthContext        string
	AuthJWTCredentials string
)

type UserCredentials struct {
	Id              *int64
	Email 			string
	Verified_email  bool
	Name 			string
	Given_name      string
	Family_name 	string
	Picture			string
	Locale 			string
}

// Authenticator checks authentication for a given request
type Authenticator interface {
	// GetToken obtains authorization token from request
	GetToken(ctx context.Context, request *http.Request) (string, error)
	// GetUserCredentials obtains user credentials by token
	GetUserCredentials(ctx context.Context, token string) (*UserCredentials, error)
	// Authorize checks credentials for permissions and returns error if request is not allowed
	Authorize(ctx context.Context, r *http.Request, credentials UserCredentials) error
}

// JwtAuthenticator contains methods to deal with JWT inside the auth middleware
type JwtAuthenticator interface {
	// GenerateToken makes jwt token string that contains user credentials
	GenerateToken(credentials UserCredentials) (string, error)
	// ParseToken obtains user credentials from token
	ParseToken(tokenStr string) (*UserCredentials, error)
}
