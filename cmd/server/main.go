package main

import (
	"anubis/app/core"
	"anubis/app/core/helpers"
	"anubis/app/core/middlewares"
	v1 "anubis/app/routes/v1"
	"anubis/tools/databace/psql"
	"anubis/tools/utils"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"time"
)

func main() {
	env := core.Env{}
	utils.NewEnv(&env)
	/**
	* ========================
	*  Setup db
	* ========================
	 */
	var validate *validator.Validate

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
		err := validate.RegisterValidation("phone", helpers.ValidatePhone)
		if err != nil {
			log.Fatalf("CRITICAL: ", "validate Er: %v\n", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		env.PgUsername,
		env.PgPassword,
		env.PgHost,
		env.PgPort,
		env.PgDatabase,
	)
	pgStore, err := psql.NewClient(ctx, 5, 3*time.Second, databaseUrl, false)
	if err != nil {
		log.Fatalf("CRITICAL: ", "unexpected error while tried to connect to database.md: %v\n", err)
	}
	defer pgStore.Close()

	/**
	* ========================
	*  Setup Application
	* ========================
	 */

	app := gin.New()
	app.Use(gin.Recovery())

	if env.AppEnv != "development" {
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

	app.Use(middlewares.RequestID(nil))
	app.Use(gin.LoggerWithConfig(core.GetLoggerConfig(nil, nil, nil)))
	/**
	* ========================
	* Initialize All Route
	* ========================
	 */
	app.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	v1.NewRouteUser(pgStore, app, env)

	protectedRouter := app.Group("/check_auth")
	protectedRouter.Use(middlewares.JwtAuthMiddleware(env.AccessTokenSecret))
	protectedRouter.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "auth"}) })

	start := app.Run(env.AppIp)
	if start != nil {
		log.Fatalf("unexpected error while tried to start localhost: %v\n", err)
	}
}
