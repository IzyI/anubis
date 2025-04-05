package middlewares

import (
	"anubis/app/core"
	"anubis/app/core/schemes"
	"anubis/tools/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	UserIDKey   = "user-id"
	UserRole    = "user-role"
	UserProject = "user-project"
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
			var authClaims schemes.JwtCustomClaims
			err := utils.ExtractToken(authToken, config.AccessTokenSecret, &authClaims)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, schemes.HTTPError{
					Code: 97,
					Err:  "BAD token",
				})
				return
			}
			data := strings.Split(authClaims.D, "|")
			if len(data) >= 2 {
				c.Set(UserProject, data[0])
				c.Set(UserRole, data[1])
			}
			c.Set(UserIDKey, authClaims.RegisteredClaims.Subject)
			c.Request.Header.Set(UserIDKey, authClaims.RegisteredClaims.Subject)
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
