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
	Service = "service"
	Domain  = "domain"
	Once    = "once"
	Auth    = "auth"
)

func DomainMiddleware(config core.ServiceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := c.Request.Header.Get("Domain")
		if domain == "" {
			host := strings.Split(c.Request.Host, ".")
			ll := len(host)
			if ll < 2 {
				domain = host[0]
			} else {
				domain = host[ll-2]
			}

		}
		lService, ok := config.ListServices[domain]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &schemes.HTTPError{Code: 104, Err: "Domain not found"})
			return
		}
		c.Set(Service, lService.Service)
		c.Set(Domain, domain)
		c.Set(Once, lService.Once)
		c.Set(Auth, lService.Auth)
		return
	}
}

func CheckOnceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(Once) {
			c.AbortWithStatusJSON(http.StatusForbidden, &schemes.HTTPError{Code: 106, Err: "Authorization method denied"})
			return
		}
		c.Next()
	}
}

func CheckAuthTypeMiddleware(typeAuth string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authIface, exists := c.Get(Auth)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, &schemes.HTTPError{Code: 106, Err: "Authorization data missing"})
			return
		}

		// Приводим к []string
		auth, ok := authIface.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, &schemes.HTTPError{Code: 106, Err: "Invalid authorization data"})
			return
		}

		if !utils.LittleContainsString(auth, typeAuth) {
			c.AbortWithStatusJSON(http.StatusForbidden, &schemes.HTTPError{Code: 106, Err: "Authorization method denied"})
			return
		}

		c.Next()
	}
}
