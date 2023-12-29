package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/bdarge/api-gateway/cmd/docs"
	"github.com/bdarge/api-gateway/pkg/auth"
	"github.com/bdarge/api-gateway/pkg/config"
	"github.com/bdarge/api-gateway/pkg/customer"
	"github.com/bdarge/api-gateway/pkg/profile"
	"github.com/bdarge/api-gateway/pkg/transaction"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/jub0bs/fcors"
	"github.com/jub0bs/fcors/risky"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/exp/slog"
	"log"
	"net/http"
	"os"
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
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "dev"
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("start app")

	conf, err := config.LoadConfig(environment)

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	// Creates a router without any middleware by default
	router := gin.New()

	if err = router.SetTrustedProxies(nil); err != nil {
		log.Fatalln("Failed at SetTrustedProxies", err)
	}

	// By default, gin.DefaultWriter = os.Stdout, change the format
	router.Use(jsonLoggerMiddleware())
	// router.Use(slog.Logger{})

	slog.Info("configure cors")

	//reading https://jub0bs.com/posts/2023-02-08-fearless-cors/#3-provide-support-for-private-network-access
	cors, corsErr := fcors.AllowAccess(
		fcors.FromOrigins("https://localhost:4201", "http://localhost:4201"),
		fcors.WithMethods(
			http.MethodGet,
			http.MethodDelete,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodOptions,
		),
		fcors.WithRequestHeaders(
			"Authorization",
			"Content-Type",
		),
		fcors.MaxAgeInSeconds(30),
		risky.PrivateNetworkAccess(),
	)

	if corsErr != nil {
		log.Fatalln("Failed at CORS setup", corsErr)
	}

	// apply the CORS middleware to the router
	router.Use(adapter.Wrap(cors))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "")
	})

	slog.Info("configure doc")

	//url := ginSwagger.URL(fmt.Sprintf("%s/docs/swagger.json", conf.BaseUrl))
	//router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	slog.Info("set routes")
	v1 := router.Group("/v1")
	{
		v1.GET("/docs", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s/docs/index.html", conf.BaseUrl))
		})
		v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		authSvc := *auth.RegisterRoutes(v1, &conf)
		transaction.RegisterRoutes(v1, &conf, &authSvc)
		customer.RegisterRoutes(v1, &conf, &authSvc)
		profile.RegisterRoutes(v1, &conf, &authSvc)
	}

	if err = router.Run(conf.Port); err != nil {
		log.Fatalln("Failed at gin.Run", err)
	}
	slog.Info("api gateway started", "port", conf.Port)
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
