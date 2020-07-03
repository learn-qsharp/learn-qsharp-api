package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/learn-qsharp/learn-qsharp-api/api"
)

func Run(db *gorm.DB) error {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/tutorials/:id", api.ShowTutorial)
	r.GET("/tutorials", api.ListTutorials)

	return r.Run()
}
