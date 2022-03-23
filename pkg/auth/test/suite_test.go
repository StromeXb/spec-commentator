package test

import (
	"context"
	"fmt"
	"testing"

	"spec-commentor/config"

	"spec-commentor/pkg/auth"

	"spec-commentor/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type JwtSuite struct {
	suite.Suite

	ctx    context.Context
	logger *zerolog.Logger
	cfg    *auth.JWTConfig
}

func (s *JwtSuite) SetupSuite() {
	var err error
	s.ctx = context.Background()

	s.logger, err = logger.NewLogger(&logger.Config{
		LogLevel: "debug",
	})
	s.dieOnErr(err)

	s.cfg = config.NewJWTConfig(config.NewViper())
}

func (s *JwtSuite) dieOnErr(err error, msg ...string) {
	if err != nil {
		s.FailNow(fmt.Sprintf("Unexpected error %s", msg), err)
	}
}

func (s *JwtSuite) TearDownSuite() {}

func TestJwtRunner(t *testing.T) {
	suite.Run(t, new(JwtSuite))
}
