package routes

import (
	"anubis/app/api/DAL/storage"
	"anubis/app/api/controllers"
	"anubis/app/api/usecase"
	"anubis/app/core"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouteUser(db *mongo.Client, router *gin.Engine, config core.ServiceConfig) {
	repositoryMongoAuth := storage.NewRepositoryMongoAuthPhone(db)
	repositoryMongoUser := storage.NewRepositoryMongoUser(db)
	repositoryMongoProject := storage.NewRepositoryMongoProject(db)
	serviceAuth := usecase.NewServiceAuth(repositoryMongoUser, repositoryMongoProject, repositoryMongoAuth, config)
	handler := controllers.NewControllerAuth(serviceAuth)

	route := router.Group("/auth")
	//
	//route.Use(middlewares.AuthToken())
	//route.Use(middlewares.AuthRole(map[string]bool{"admin": true, "merchant": true}))
	//
	route.POST("/reg", handler.HandlerRegPOST)
	route.POST("/valid", handler.HandlerValidSmsPOST)
	route.POST("/login_phone", handler.HandlerLoginPOST)
	//route.POST("/refresh", handler.HandlerRefreshTokenPOST)

	//route.GET("/result/:id", handler.HandlerResult)
	//route.DELETE("/delete/:id", handler.HandlerDelete)
	//route.PUT("/update:id", handler.HandlerUpdate)
}
