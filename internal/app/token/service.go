package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	signingKey []byte
}

func New(signingKey string) *Service {
	return &Service{
		signingKey: []byte(signingKey),
	}
}

func (s *Service) Generate(userID uint64) (string, error) {
	payload := jwt.MapClaims{
		"sub": userID,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS512, payload).SignedString(s.signingKey)
}

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)

func (s *Service) Validate(token string) (uint64, error) {
	claims := jwt.MapClaims{}

	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return s.signingKey, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, ErrInvalidToken
	}

	return uint64(claims["sub"].(float64)), nil
}
