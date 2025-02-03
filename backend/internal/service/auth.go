package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nessai1/aiinterview/internal/domain"
	"time"
)

type AuthService struct {
	secret string
}

type claims struct {
	jwt.RegisteredClaims
	UserUUID string
}

const tokenExp = time.Hour * 24 * 31

func (s *AuthService) FetchUserFromToken(tokenString string) (domain.User, error) {
	c := claims{}
	token, err := jwt.ParseWithClaims(tokenString, &c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: token must have HMAC signing algorithm")
		}

		return []byte(s.secret), nil
	})

	if err != nil {
		return domain.User{}, fmt.Errorf("cannot parse JWT token: %w", err)
	}

	if !token.Valid {
		return domain.User{}, fmt.Errorf("invalid JWT token")
	}

	return domain.User{UUID: c.UserUUID}, nil
}

func (s *AuthService) BuildTokenByUser(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserUUID: user.UUID,
	})

	tokenStr, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("cannot sign token: %w", err)
	}

	return tokenStr, nil
}
