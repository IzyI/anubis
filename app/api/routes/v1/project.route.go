package routes

import (
	"anubis/app/api/DAL/storage"
	controllers2 "anubis/app/api/api/controllers"
	usecase2 "anubis/app/api/api/usecase"
	"anubis/app/core"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouteProject(db *mongo.Client, app *gin.RouterGroup, config core.ServiceConfig) {
	repositoryMongoProject := storage.NewRepositoryMongoProject(db)
	repositoryMongoUser := storage.NewRepositoryMongoUser(db)

	serviceProject := usecase2.NewServiceProject(repositoryMongoUser, repositoryMongoProject, config)
	handler := controllers2.NewControllerProjects(serviceProject)

	serviceProjectMembers := usecase2.NewServiceServiceProjectMembers(repositoryMongoUser, repositoryMongoProject, config)
	handler2 := controllers2.NewControllerProjectMembers(serviceProjectMembers)

	route := app.Group("")

	route.GET("", handler.HandlerGetProjectsFlow)
	route.POST("", handler.HandlerPostProjectsFlow)
	route.DELETE("/:ID", handler.HandlerDelProjectsFlow)
	route.GET("/:ID", handler2.HandlerGetProjectMembersFlow)
	route.POST("/:ID/member", handler2.HandlerPOSTProjectMembersFlow)
	route.PUT("/:ID/member", handler2.HandlerPUTProjectMembersFlow)
	route.DELETE("/:ID/member", handler2.HandlerDELProjectMembersFlow)
}
