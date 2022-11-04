package main

import (
	_ "github.com/bdarge/sb-api-gateway/cmd/docs"
	"github.com/bdarge/sb-api-gateway/pkg/auth"
	"github.com/bdarge/sb-api-gateway/pkg/config"
	"github.com/bdarge/sb-api-gateway/pkg/disposition"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

// @title SM Swagger API
// @version 1.0
// @description Swagger API for Business X.
// @termsOfService http://swagger.io/terms/

// @BasePath /v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	conf, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	// Creates a router without any middleware by default
	router := gin.New()

	//Default() allows all origins
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default, gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	url := ginSwagger.URL("http://127.0.0.1:3000/swagger/doc.json")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	v1 := router.Group("/v1")
	{
		authSvc := *auth.RegisterRoutes(v1, &conf)
		disposition.RegisterRoutes(v1, &conf, &authSvc)
	}

	router.Run(conf.Port)
}
