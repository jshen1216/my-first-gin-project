package main

import (
	"golangAPI/src"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	finder := server.Group("/accesslog")
	src.FindDataRouter(finder)
	server.Run(":8080")
}
