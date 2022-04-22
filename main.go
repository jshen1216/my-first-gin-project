package main

import (
	"golangAPI/src"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	mySQL := server.Group("/mysql")
	src.FindDataRouter(mySQL)
	server.Run(":8080")
}
