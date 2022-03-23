package test

import (
	"time"

	"spec-commentor/pkg/auth"
	"spec-commentor/pkg/clock"
)

// кодирование JWT
func (s *JwtSuite) TestJwtEncode() {
	cl := clock.NewMock()
	cl.On("Now").Return(time.Parse("2006-01-02T15:04:05+07:00", "2006-01-02T15:04:05+07:00"))

	authenticator, err := auth.NewJwtAuthenticator(
		&auth.JWTConfig{
			PublicKeyPath:      s.cfg.PublicKeyPath,
			PrivateKeyPath:     s.cfg.PrivateKeyPath,
			JwtLifeTimeMinutes: 15,
		},
		cl,
	)
	s.dieOnErr(err)

	userID := int64(1)
	userCredentials := auth.UserCredentials{
		UserID: &userID,
		Role:   auth.UserRoleAdmin,
	}

	token, err := authenticator.GenerateToken(userCredentials)
	s.dieOnErr(err)

	s.Equal(token, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsIlJvbGUiOiJhZG1pbiIsIkF1dGhCYWNrQWNjZXNzVG9rZW4iOiIiLCJleHAiOjExMzYyMTUxNDUsImlhdCI6MTEzNjIxNDI0NX0.sFg7FJVW36WEp0_WH6lm81n-KoDnTprtjM8yYKvjwahQz5LhfQ5iV_VwCFycv-0Z36bcG4n4gtCKlpo0IrpHc-TfLjhYY3Z4PqgneNy--bxzXcLJUVYr-nBesEF4ZQ7fHCUbrg_q4_Qcax5pdicelL5aTAn8cWNt8J2fpmvTzRd9IPPqMzB-JUQnrt3q7aXuN9JKv9PjEVaQwUHHlq7StgvCkEmln_KNFFX1Vu-Z2YnJG68wKQdk0ZfhsNNoAGGYEoKqBdQEqm2GCWjmItPQgKpHG768KyejOm6eu1o2MQwWWim88oe8nfQfnWgu1FOPLWy2uHzfXALPcq_ynw2Qjgqn3GJzb6Bi4oEgveZAz7kAs33Z7LyHF_tl38rY9fnQrxvaxtZOYWfBCTOe4-tGbvjDM1ngaY39yjyjUx0OjmfNjkgoDLmyZvLYVkbFhLO3yuwXfbNdH6aFsCYbeESkT6qn5_MBOXUVya8MS8uX_4jOGZXMG7LykVY_ESNcWLqucSUcz3a2rRn8nO3KY6ZvcbIjLV3NLLDsJy-1u5RR5cSm3Mr2caaDoZQayZTY7NVm0OSFfwI8B24SdYzSP9hbJmCI5L8hBMhwTp3dsGpT5KnKDXwvRh531zvVXao5XvEoZYqYWVrSSvxMKn0vuXbpescRqkD39v351mpybovTXEk")
}

// раскодирование JWT
func (s *JwtSuite) TestJwtDecode() {
	userID := int64(1)

	cred := auth.UserCredentials{
		UserID: &userID,
		Role:   auth.UserRoleAdmin,
	}

	authenticator, err := auth.NewJwtAuthenticator(
		&auth.JWTConfig{
			PublicKeyPath:      s.cfg.PublicKeyPath,
			PrivateKeyPath:     s.cfg.PrivateKeyPath,
			JwtLifeTimeMinutes: 15,
		},
		clock.New(),
	)
	s.dieOnErr(err)

	token, err := authenticator.GenerateToken(cred)
	s.dieOnErr(err)

	newCred, err := authenticator.ParseToken(token)
	s.dieOnErr(err)

	s.Equal(newCred, &cred)
}

func (s *JwtSuite) TestJwtDecodeWithoutPrivateToken() {
	userID := int64(1)
	cred := auth.UserCredentials{
		UserID: &userID,
		Role:   auth.UserRoleAdmin,
	}

	authenticator, err := auth.NewJwtAuthenticator(
		&auth.JWTConfig{
			PublicKeyPath:      s.cfg.PublicKeyPath,
			PrivateKeyPath:     s.cfg.PrivateKeyPath,
			JwtLifeTimeMinutes: 15,
		},
		clock.New(),
	)
	s.dieOnErr(err)

	token, err := authenticator.GenerateToken(cred)
	s.dieOnErr(err)

	authenticatorDecode, err := auth.NewJwtAuthenticator(
		&auth.JWTConfig{
			PublicKeyPath:      s.cfg.PublicKeyPath,
			PrivateKeyPath:     "",
			JwtLifeTimeMinutes: 15,
		},
		clock.New(),
	)
	s.dieOnErr(err)

	newCred, err := authenticatorDecode.ParseToken(token)
	s.dieOnErr(err)

	s.Equal(newCred, &cred)
}
