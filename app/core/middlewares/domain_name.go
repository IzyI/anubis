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
	Service = "name_service"
	Domain  = "name_domain"
)

func CheckDomain(s core.ServiceConfig, d string) (string, error) {
	println(d)
	domain, ok := s.ListServices[d]
	if !ok {
		return "", &schemes.HTTPError{Code: 104, Err: "Domain not found"}
	}
	if domain.Auth != nil {
		if !utils.LittleContainsString(domain.Auth, "phone") {
			return "", &schemes.HTTPError{Code: 106, Err: "Authorization method denied"}
		}
	} else {
		return "", &schemes.HTTPError{Code: 106, Err: "Authorization method denied !"}
	}
	return domain.Service, nil
}

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
		service, err := CheckDomain(config, domain)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}
		c.Set(Service, service)
		c.Set(Domain, domain)
		c.Next()
		return
	}
}
