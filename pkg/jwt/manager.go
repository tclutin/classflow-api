package jwt

import (
	"errors"
	jsontoken "github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type Manager interface {
	NewToken(userID uint64, ttl time.Duration) (string, error)
	ParseToken(accessToken string) (uint64, error)
}

type TokenManager struct {
	signingKey string
}

func MustLoadTokenManager(signingKey string) *TokenManager {
	if signingKey == "" {
		log.Fatalln("signingKey is empty")
	}

	return &TokenManager{signingKey: signingKey}
}

func (t *TokenManager) NewToken(userID uint64, ttl time.Duration) (string, error) {
	payload := jsontoken.MapClaims{
		"exp": time.Now().Add(ttl).Unix(),
		"sub": userID,
	}

	token := jsontoken.NewWithClaims(jsontoken.SigningMethodHS256, payload)

	return token.SignedString([]byte(t.signingKey))
}

func (t *TokenManager) ParseToken(jwtToken string) (uint64, error) {
	token, err := jsontoken.Parse(jwtToken, func(token *jsontoken.Token) (interface{}, error) {
		if _, ok := token.Method.(*jsontoken.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(t.signingKey), nil
	})

	if err != nil {
		return 0, errors.New("token is expired")
	}

	claims, ok := token.Claims.(jsontoken.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims format")
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, errors.New("sub claim missing")
	}

	subFloat, ok := sub.(float64)
	if !ok {
		return 0, errors.New("invalid sub format, expected float64 or string")
	}

	return uint64(subFloat), nil
}
