package routes

import (
	storage2 "anubis/app/DAL/storage"
	"anubis/app/api/controllers"
	"anubis/app/api/usecase"
	"anubis/app/core"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouteProject(db *mongo.Client, app *gin.RouterGroup, config core.ServiceConfig) {
	repositoryMongoProject := storage2.NewRepositoryMongoProject(db, &config)
	repositoryMongoUser := storage2.NewRepositoryMongoUser(db)

	serviceProject := usecase.NewServiceProject(repositoryMongoUser, repositoryMongoProject, config)
	handler := controllers.NewControllerProjects(serviceProject)

	serviceProjectMembers := usecase.NewServiceServiceProjectMembers(repositoryMongoUser, repositoryMongoProject, config)
	handler2 := controllers.NewControllerProjectMembers(serviceProjectMembers)

	route := app.Group("")

	route.GET("", handler.HandlerGETProjects)
	route.POST("", handler.HandlerPOSTProjects)
	route.DELETE("/:id", handler.HandlerDELProjects)
	route.PUT("/:id", handler.HandlerPUTProjects)
	route.GET("/:id", handler2.HandlerGETProjectMembers)
	route.POST("/:id/member", handler2.HandlerPOSTProjectMembers)
	route.PUT("/:id/member", handler2.HandlerPUTProjectMembers)
	route.DELETE("/:id/member", handler2.HandlerDELProjectMembers)
}
