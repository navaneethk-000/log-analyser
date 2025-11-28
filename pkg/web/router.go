package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupRouter(db *gorm.DB) *gin.Engine {
	DB = db

	r := gin.Default()
	r.Use(PrintHelloWorldBeforeRouting)
	r.Use(CORSMiddleware())

	//allow domain = "*"
	r.LoadHTMLGlob("cmd/web/templates/*")

	// r.GET("/", GetAllLogs)
	// r.POST("/filter", ExecuteFilterQuery)
	// r.POST("/", PaginatedLogs)
	r.POST("/filter", FilterPaginatedLogs)

	return r
}

func PrintHelloWorldBeforeRouting(ctx *gin.Context) {
	fmt.Println("Hello World")
	ctx.Next()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
