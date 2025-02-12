package utils

import (
	"anubis/app/core/schemes"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

func CreateAccessToken(uuid string, group []string, secret string, expiry int) (accessToken string, err error) {
	claims := &schemes.JwtCustomClaims{
		ID: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiry))),
			Audience:  group,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(uuid string, group []string, secret string, expiry int) (refreshToken string, err error) {
	claims := &schemes.JwtCustomRefreshClaims{
		ID: uuid,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiry))),
			Audience:  group,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("invalid Token")
	}

	return claims["id"].(string), nil
}

func RemoveFirstPart(token string) string {
	if i := strings.IndexByte(token, '.'); i >= 0 {
		return token[i+1:]
	}
	return token // если точка не найдена
}

const FIRSTPARTJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
