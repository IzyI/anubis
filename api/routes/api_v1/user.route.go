package api_v1

import (
	"anubis/api/controllers"
	"anubis/api/repositories"
	"anubis/core"
	"anubis/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouteUser(db *pgxpool.Pool, router *gin.Engine, env core.Env) {
	repository := repositories.NewRepositoryUser(db)
	service := services.NewServiceUser(repository, env)
	handler := controllers.NewControllerUser(service)

	route := router.Group("/auth")
	//
	//route.Use(middlewares.AuthToken())
	//route.Use(middlewares.AuthRole(map[string]bool{"admin": true, "merchant": true}))
	//
	route.POST("/reg", handler.HandlerRegPOST)
	route.POST("/valid", handler.HandlerValidSmsPOST)
	route.POST("/login", handler.HandlerLoginPOST)
	route.POST("/refresh", handler.HandlerRefreshTokenPOST)
	//route.GET("/result/:id", handler.HandlerResult)
	//route.DELETE("/delete/:id", handler.HandlerDelete)
	//route.PUT("/update:id", handler.HandlerUpdate)
}
