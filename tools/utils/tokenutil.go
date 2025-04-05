package utils

import (
	"anubis/app/core/schemes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

func CreateAccessToken(id string, data string, domains []string, secret string, expiry int) (accessToken string, err error) {
	claims := &schemes.JwtCustomClaims{
		D: data,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiry))),
			Audience:  domains,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(idSession string, secret string, expiry int) (refreshToken string, err error) {
	claims := &schemes.JwtCustomRefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiry))),
			ID:        idSession,
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

func ExtractToken[T jwt.Claims](requestToken string, secret string, claims T) error {
	token, err := jwt.ParseWithClaims(
		requestToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)

	if err != nil {
		return fmt.Errorf("token parsing failed: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func RemoveFirstPart(token string) string {
	if i := strings.IndexByte(token, '.'); i >= 0 {
		return token[i+1:]
	}
	return token // если точка не найдена
}

func ParseInfoToken(input string) (string, string, string, error) {
	// Разбиваем строку по разделителю "|"
	parts := strings.Split(input, "|")

	// Проверяем, что у нас ровно 3 части (userID, projectID, роль)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("неверный формат входной строки")
	}

	userID := parts[0]
	projectID := parts[1]
	role := parts[2]

	return userID, projectID, role, nil
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	//ID string `json:"id"`
	//Roles []string `json:"roles"`
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
