package auth

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"
	"time"

	"spec-commentor/pkg/test"

	"spec-commentor/pkg/clock"
	"github.com/dgrijalva/jwt-go"
)

const (
	JwtPublicKeyPath   = "/data/jwt/cred.pub"
	JwtPrivateKeyPath  = "/data/jwt/cred"
	JwtSecret          = "secret"
	JwtLifeTimeMinutes = 1209600
)

type customUserClaims struct {
	UserCredentials
	jwt.StandardClaims
}

type jwtAuthenticator struct {
	publicKey            []byte
	privateKey           []byte
	tokenLifetimeMinutes int
	clock                clock.Clock
}

// Generate RSA key-pair
//
// openssl genrsa -out cert/id_rsa 4096
//
// openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub
type JWTConfig struct {
	PublicKeyPath      string
	PrivateKeyPath     string
	JwtLifeTimeMinutes int
}

func NewJwtAuthenticator(cfg *JWTConfig, cl clock.Clock) (JwtAuthenticator, error) {
	var err error
	var prvKey []byte
	var pubKey []byte
	var privateKeyPath string
	var publicKeyPath string

	if test.IsTest() {
		dir, err := build.Default.Import("git.redmadrobot.com/internship/backend/lim-ext/", "", build.FindOnly)
		if err != nil {
			return nil, err
		}
		privateKeyPath = filepath.Join(dir.Dir, cfg.PrivateKeyPath)
		publicKeyPath = filepath.Join(dir.Dir, cfg.PublicKeyPath)
	} else {
		privateKeyPath = cfg.PrivateKeyPath
		publicKeyPath = cfg.PublicKeyPath
	}

	if cfg.PrivateKeyPath != "" {
		prvKey, err = ioutil.ReadFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read private key: %s", err)
		}
	}

	pubKey, err = ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key: %s", err)
	}

	return &jwtAuthenticator{
		publicKey:            pubKey,
		privateKey:           prvKey,
		tokenLifetimeMinutes: cfg.JwtLifeTimeMinutes,
		clock:                cl,
	}, nil
}

func (j *jwtAuthenticator) GenerateToken(credentials UserCredentials) (string, error) {
	uc := customUserClaims{
		credentials,
		jwt.StandardClaims{
			IssuedAt:  j.clock.Now().Unix(),
			ExpiresAt: j.clock.Now().Add(time.Minute * time.Duration(j.tokenLifetimeMinutes)).Unix(),
		},
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("generate tokenn: parse key: %w", err)
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, uc).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *jwtAuthenticator) ParseToken(tokenStr string) (*UserCredentials, error) {
	uc := customUserClaims{}

	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return nil, fmt.Errorf("validate token: parse key: %w", err)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*customUserClaims); ok && token.Valid {
		return &claims.UserCredentials, nil
	}

	return nil, err
}
