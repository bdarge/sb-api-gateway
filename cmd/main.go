package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/bdarge/api-gateway/cmd/docs"
	"github.com/bdarge/api-gateway/pkg/auth"
	"github.com/bdarge/api-gateway/pkg/config"
	"github.com/bdarge/api-gateway/pkg/currency"
	"github.com/bdarge/api-gateway/pkg/customer"
	"github.com/bdarge/api-gateway/pkg/lang"
	"github.com/bdarge/api-gateway/pkg/profile"
	"github.com/bdarge/api-gateway/pkg/transaction"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/exp/slog"
)

//	@title			SM Swagger API
//	@version		1.0
//	@description	Swagger API for Business X.
//	@termsOfService	http://swagger.io/terms/

//	@BasePath	/v1

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
func main() {
	var programLevel = new(slog.LevelVar)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
	slog.SetDefault(logger)

	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "dev"
	}

	conf, err := config.LoadConfig(environment)

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	programLevel.Set(conf.LogLevel)

	slog.Info("Start api-gateway")

	// Creates a router without any middleware by default
	router := gin.New()

	if err = router.SetTrustedProxies(nil); err != nil {
		log.Fatalln("Failed at SetTrustedProxies", err)
	}

	// By default, gin.DefaultWriter = os.Stdout, change the format
	router.Use(jsonLoggerMiddleware())

	slog.Info("Configure cors")

	c := cors.New(cors.Options{
    AllowedOrigins: []string{"http://localhost:*", "http://sb.odainfo.com"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
    // Enable Debugging for testing, consider disabling in production
    Debug: false,
	})

	// apply the CORS middleware to the router
	router.Use(c)

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})

	slog.Info("set routes")
	v1 := router.Group("/v1")
	{
		slog.Info("configure doc")
		v1.GET("/docs", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s/docs/index.html", conf.BaseURL))
		})
		v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		authSvc := *auth.RegisterRoutes(v1, &conf)
		transaction.RegisterRoutes(v1, &conf, &authSvc)
		customer.RegisterRoutes(v1, &conf, &authSvc)
		profile.RegisterRoutes(v1, &conf, &authSvc)
		currency.RegisterRoutes(v1, &conf, &authSvc)
		lang.RegisterRoutes(v1, &conf, &authSvc)
	}

	if err = router.Run(conf.Port); err != nil {
		log.Fatalln("Failed at gin.Run", err)
	}
	slog.Info("Gateway started", "port", conf.Port)
}

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			glog := make(map[string]interface{})

			glog["status_code"] = params.StatusCode
			glog["path"] = params.Path
			glog["method"] = params.Method
			glog["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			glog["remote_addr"] = params.ClientIP
			glog["response_time"] = params.Latency.String()

			s, _ := json.Marshal(glog)
			return string(s) + "\n"
		},
	)
}
