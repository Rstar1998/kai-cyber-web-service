package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/scan", ScanRepo) // Scan POST API
	r.POST("/query", QueryVulnerabilities) // Query POST API
	return r
}

func main() {
	InitDB()  // Initilize our sqlite DB
	r := setupRouter()
	r.Run(":8080") // Run on port 8080
}
