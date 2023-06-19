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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("start app...")

	conf, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	// Creates a router without any middleware by default
	router := gin.New()

	if err = router.SetTrustedProxies(nil); err != nil {
		log.Fatalln("Failed at SetTrustedProxies", err)
	}

	//Default() allows all origins
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	// By default, gin.DefaultWriter = os.Stdout, change the format
	router.Use(jsonLoggerMiddleware())
	// router.Use(slog.Logger{})

	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	url := ginSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", conf.DocUrl))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	v1 := router.Group("/v1")
	{
		authSvc := *auth.RegisterRoutes(v1, &conf)
		transaction.RegisterRoutes(v1, &conf, &authSvc)
		customer.RegisterRoutes(v1, &conf, &authSvc)
		profile.RegisterRoutes(v1, &conf, &authSvc)
	}

	if err = router.Run(conf.Port); err != nil {
		log.Fatalln("Failed at gin.Run", err)
	}
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
