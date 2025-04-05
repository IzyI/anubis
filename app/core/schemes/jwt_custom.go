package schemes

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	D string `json:"D"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	//ID string `json:"id"`
	jwt.RegisteredClaims
}
