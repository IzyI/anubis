package routes

import (
	storage2 "anubis/app/DAL/storage"
	"anubis/app/api/controllers"
	"anubis/app/api/usecase"
	"anubis/app/core"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouteToken(db *mongo.Client, app *gin.RouterGroup, config core.ServiceConfig) {
	repositoryMongoProject := storage2.NewRepositoryMongoProject(db, &config)
	repositoryMongoUser := storage2.NewRepositoryMongoUser(db)
	serviceToken := usecase.NewServiceToken(repositoryMongoProject, repositoryMongoUser, config)
	handler := controllers.NewControllerToken(serviceToken)

	route := app.Group("")
	route.POST("/refresh_token", handler.HandlerPOSTRefreshTokenDomain)
	route.POST("/logout", handler.HandlerPOSTLogoutDomain)
}
