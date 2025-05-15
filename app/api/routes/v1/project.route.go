package routes

import (
	storage2 "anubis/app/DAL/storage"
	"anubis/app/api/controllers"
	"anubis/app/api/usecase"
	"anubis/app/core"
	"anubis/app/core/middlewares"
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
	//TODO: у меня при создании проджекта есть возможность создать проджекты с одинаковым именем
	route.POST("", middlewares.CheckOnceMiddleware(), handler.HandlerPOSTProjects)
	route.DELETE("/:id", middlewares.CheckOnceMiddleware(), handler.HandlerDELProjects)
	route.PUT("/:id", middlewares.CheckOnceMiddleware(), handler.HandlerPUTProjects)
	route.GET("/:id", middlewares.CheckOnceMiddleware(), handler2.HandlerGETProjectMembers)
	route.POST("/:id/member", middlewares.CheckOnceMiddleware(), handler2.HandlerPOSTProjectMembers)
	route.PUT("/:id/member", middlewares.CheckOnceMiddleware(), handler2.HandlerPUTProjectMembers)
	route.DELETE("/:id/member", middlewares.CheckOnceMiddleware(), handler2.HandlerDELProjectMembers)
}
