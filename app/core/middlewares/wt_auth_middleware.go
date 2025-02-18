package middlewares

import (
	"anubis/app/core"
	"anubis/app/core/schemes"
	"anubis/tools/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	UserIDKey = "user-id"
)

func JwtAuthMiddleware(config core.ServiceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authToken string

		// First, check if the token is present in cookies
		if cookie, err := c.Cookie("Token"); err == nil {
			authToken = cookie
		} else {
			authToken = c.Request.Header.Get("Token")
			if len(authToken) == 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, schemes.HTTPError{
					Code: 99,
					Err:  "Unable to find token",
				})
				return
			}
		}

		if config.ShortJwt {
			authToken = config.ShortJwtValue + "." + authToken
		}
		// Validate the token
		authorized, _ := utils.IsAuthorized(authToken, config.AccessTokenSecret)
		if authorized {
			var refreshClaims utils.RefreshClaims
			err := utils.ExtractToken(authToken, config.AccessTokenSecret, &refreshClaims)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, schemes.HTTPError{
					Code: 97,
					Err:  "User not found",
				})
				c.Abort()
				return
			}
			c.Set(UserIDKey, refreshClaims.ID)
			c.Request.Header.Set(UserIDKey, refreshClaims.ID)
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, schemes.HTTPError{
			Code: 98,
			Err:  "Not authorized",
		})
		return
	}
}
