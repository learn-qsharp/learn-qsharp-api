package router

import "github.com/gin-gonic/gin"

func Run() error {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	return r.Run()
}
