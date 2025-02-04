package routes

import (
	"anubis/app/api/controllers"
	"anubis/app/api/storage"
	"anubis/app/api/usecase"
	"anubis/app/core"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouteUser(db *pgxpool.Pool, router *gin.Engine, config core.ServiceConfig) {
	repositoryPsqlAuth := storage.NewRepositoryPsqlAuthPhone(db)
	repositoryPsqlUser := storage.NewRepositoryPsqlUser(db)
	serviceAuth := usecase.NewServiceAuth(repositoryPsqlAuth, repositoryPsqlUser, config)
	handler := controllers.NewControllerAuth(serviceAuth)

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
