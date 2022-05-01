package main

import (
	"golangAPI/src"

	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "golangAPI/docs"
)

// @title           Swagger API
// @version         1.0
// @description     This is sample of practicing API to get data in MySQL & ES
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:8080
// @BasePath  /accesslog

func main() {
	server := gin.Default()
	finder := server.Group("/accesslog")
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	src.FindDataRouter(finder)
	server.Run(":8080")
}
