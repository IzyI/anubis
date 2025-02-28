package routes

import (
	"anubis/app/api/DAL/storage"
	"anubis/app/api/api/controllers"
	"anubis/app/api/api/usecase"
	"anubis/app/core"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRoutePhoneAuth(db *mongo.Client, app *gin.RouterGroup, config core.ServiceConfig) {
	repositoryMongoAuth := storage.NewRepositoryMongoAuthPhone(db)
	repositoryMongoUser := storage.NewRepositoryMongoUser(db)
	repositoryMongoProject := storage.NewRepositoryMongoProject(db)
	serviceAuth := usecase.NewServiceAuth(repositoryMongoUser, repositoryMongoProject, repositoryMongoAuth, config)
	handler := controllers.NewControllerAuth(serviceAuth)

	route := app.Group("/phone")
	//
	//route.Use(middlewares.AuthToken())
	//route.Use(middlewares.AuthRole(map[string]bool{"admin": true, "merchant": true}))
	//
	route.POST("/reg", handler.HandlerRegPOST)
	route.POST("/sms", handler.HandlerValidSmsPOST)
	route.POST("/login", handler.HandlerLoginPOST)
	//route.POST("/refresh", handler.HandlerRefreshTokenPOST)

	//route.GET("/result/:id", handler.HandlerResult)
	//route.DELETE("/delete/:id", handler.HandlerDelete)
	//route.PUT("/update:id", handler.HandlerUpdate)
}
