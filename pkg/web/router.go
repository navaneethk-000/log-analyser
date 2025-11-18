package web

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupRouter(db *gorm.DB) *gin.Engine {
	DB = db

	r := gin.Default()

	r.LoadHTMLGlob("cmd/web/templates/*")

	r.GET("/", GetAllLogs)
	r.POST("/filter", ExecuteFilterQuery)

	return r
}
