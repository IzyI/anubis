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

func NewRoutePhoneAuth(db *mongo.Client, app *gin.RouterGroup, config core.ServiceConfig) {
	repositoryMongoAuth := storage2.NewRepositoryMongoAuthPhone(db)
	repositoryMongoUser := storage2.NewRepositoryMongoUser(db)
	repositoryMongoProject := storage2.NewRepositoryMongoProject(db, &config)
	serviceAuth := usecase.NewServicePhoneAuth(repositoryMongoUser, repositoryMongoProject, repositoryMongoAuth, config)
	handler := controllers.NewControllerPhoneAuth(serviceAuth)

	route := app.Group("/phone")
	route.Use(middlewares.CheckAuthTypeMiddleware("phone"))
	route.POST("/reg", handler.HandlerPOSTReg)
	route.POST("/sms", handler.HandlerPOSTValidSms)
	route.POST("/login", handler.HandlerPOSTPhoneLogin)

}
