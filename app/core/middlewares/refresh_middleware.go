package middlewares

import (
	"anubis/app/core"
	"anubis/app/core/schemes"
	"anubis/tools/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SessionIDKey = "session-id"
	HashToken    = "hash-token"
)

func RefreshAuthMiddleware(config core.ServiceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var refreshToken string

		// First, check if the token is present in cookies
		if cookie, err := c.Cookie("Authorization"); err == nil {
			refreshToken = cookie
		} else {
			refreshToken = c.Request.Header.Get("Authorization")
			if len(refreshToken) == 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, schemes.HTTPError{
					Code: 99,
					Err:  "Unable to find token",
				})
				return
			}
		}

		// Validate the token
		authorized, _ := utils.IsAuthorized(refreshToken, config.RefreshTokenSecret)
		if authorized {
			var token schemes.JwtCustomRefreshClaims
			err := utils.ExtractToken(refreshToken, config.RefreshTokenSecret, &token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, schemes.HTTPError{
					Code: 97,
					Err:  "BAD token",
				})
				return
			}
			c.Set(HashToken, utils.HashToken(refreshToken))
			c.Set(SessionIDKey, token.ID)
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, schemes.HTTPError{
			Code: 98,
			Err:  "Not(R) authorized",
		})
		return
	}
}
