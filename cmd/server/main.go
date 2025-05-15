package main

import (
	"anubis/app/api/routes/v1"
	"anubis/app/core"
	"anubis/app/core/common"
	"anubis/app/core/middlewares"
	dtb "anubis/tools/dtb/mongo"
	"context"
	"fmt"
	"github.com/gin-contrib/gzip"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	config := core.ServiceConfig{}
	config.ReadConfig("config_server.env", "config_server.yaml")
	fmt.Printf("\nCONFIG: %+v\n\n", config)
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
		err = validate.RegisterValidation("safe_text", common.ValidateSafeText)
		if err != nil {
			log.Fatalf("CRITICAL: validate error: %v", err)
		}
		err = validate.RegisterValidation("object_id", common.ValidateObjectId)
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
		log.Fatalf("CRITICAL: ", "unexpected error while tried to connect to database_pg.md: %v\n", err)
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

	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE", "PATCH", "PUT", "OPTIONS"},
		AllowHeaders:    []string{"Content-Type", "Authorization", "Accept-Encoding"},
	}))

	app.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })

	app.Use(gin.LoggerWithConfig(common.GetLoggerConfig(nil, nil, nil)))
	app.Use(middlewares.DomainMiddleware(config))
	app.Use(gzip.Gzip(gzip.BestCompression))

	// ROUTS -----------------------------------------------------------------------------------------------------------

	auth := app.Group("/auth")
	routes.NewRoutePhoneAuth(database, auth, config)
	routes.NewRouteEmailAuth(database, auth, config)

	auth.Use(middlewares.RefreshAuthMiddleware(config))
	routes.NewRouteToken(database, auth, config)

	project := app.Group("/project")
	project.Use(middlewares.JwtAuthMiddleware(config))
	routes.NewRouteProject(database, project, config)

	protectedRouter := app.Group("/ping_auth")
	protectedRouter.Use(middlewares.JwtAuthMiddleware(config))
	protectedRouter.GET("", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	// -----------------------------------------------------------------------------------------------------------------
	start := app.Run(config.AppIp)
	if start != nil {
		log.Fatalf("unexpected error while tried to start localhost: %v\n", start)
	}
}
