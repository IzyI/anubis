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

func NewRouteEmailAuth(db *mongo.Client, app *gin.RouterGroup, config core.ServiceConfig) {
	repositoryMongoAuth := storage2.NewRepositoryMongoAuthEmail(db)
	repositoryMongoUser := storage2.NewRepositoryMongoUser(db)
	repositoryMongoProject := storage2.NewRepositoryMongoProject(db, &config)
	serviceAuth := usecase.NewServiceEmailAuth(repositoryMongoUser, repositoryMongoProject, repositoryMongoAuth, config)
	handler := controllers.NewControllerEmailAuth(serviceAuth)

	route := app.Group("/email")
	route.Use(middlewares.CheckAuthTypeMiddleware("email"))
	route.POST("/reg", handler.HandlerPOSTReg)
	route.POST("/code", handler.HandlerPOSTValidEmailCode)
	route.POST("/login", handler.HandlerPOSTEmailLogin)
}
