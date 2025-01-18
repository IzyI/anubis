package middlewares

import (
	"anubis/app/core/schemes"
	"anubis/tools/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, _ := utils.IsAuthorized(authToken, secret)
			if authorized {
				var userID, err = utils.ExtractToken(authToken, secret)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, schemes.ErrorResponse{
						Code: 97,
						Err:  "Not find User",
					})
					c.Abort()
					return
				}
				c.Set("x-user-id", userID)
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, schemes.ErrorResponse{
				Code: 98,
				Err:  "Not authorized",
			})
			c.Abort()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, schemes.ErrorResponse{
			Code: 99,
			Err:  "Unable to find token in header",
		})
		c.Abort()
	}
}
