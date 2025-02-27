package main

import (
	"anubis/app/core"
	"anubis/app/core/common"
	"anubis/app/core/middlewares"
	v1 "anubis/app/routes/v1"
	dtb "anubis/tools/dtb/mongo"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	config := core.ServiceConfig{}
	config.ReadConfig("config/.env", "config/server.yaml")
	fmt.Printf("\nCONFIG: %+v\n\n", config.ListServices)
	/**
	* ========================
	*  Setup db
	* ========================
	 */
	var validate *validator.Validate

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
		err := validate.RegisterValidation("phone", common.ValidatePhone)
		if err != nil {
			log.Fatalf("CRITICAL: validate error: %v", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := dtb.NewClientMongo(
		ctx,
		5,
		5*time.Second,
		fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/%s?authSource=admin&retryWrites=true&w=majority",
			config.MoUsername,
			url.QueryEscape(config.MoPassword),
			config.MoHost,
			config.MoPort,
			config.NameApp), 100)
	if err != nil {
		log.Fatalf("CRITICAL: ", "unexpected error while tried to connect to database.md: %v\n", err)
	}
	defer dtb.CloseMongoDBConnection(database)

	/**
	* ========================
	*  Setup Application
	* ========================
	 */

	app := gin.New()
	app.Use(gin.Recovery())

	if config.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize all middleware here
	app.Use(gzip.Gzip(gzip.BestCompression))
	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE", "PATCH", "PUT", "OPTIONS"},
		AllowHeaders:    []string{"Content-Type", "Authorization", "Accept-Encoding"},
	}))

	app.Use(gin.LoggerWithConfig(common.GetLoggerConfig(nil, nil, nil)))
	/**
	* ========================
	* Initialize All Route
	* ========================
	 */

	// Первая группа маршрутов для версии 1

	app.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	auth := app.Group("/auth")
	v1.NewRoutePhoneAuth(database, auth, config)
	v1.NewRouteToken(database, auth, config)

	protectedRouter := app.Group("/check_auth")
	protectedRouter.Use(middlewares.JwtAuthMiddleware(config))
	protectedRouter.GET("", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "auth"}) })

	start := app.Run(config.AppIp)
	if start != nil {
		log.Fatalf("unexpected error while tried to start localhost: %v\n", start)
	}
}
