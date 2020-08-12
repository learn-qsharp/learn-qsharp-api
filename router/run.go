package router

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/learn-qsharp/learn-qsharp-api/handlers"
)

func Run(db *pgxpool.Pool) error {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("ctx", context.Background())
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/tutorials/:id", handlers.ShowTutorial)
	r.GET("/tutorials", handlers.ListTutorials)

	r.GET("/problems/:id", handlers.ShowProblem)
	r.GET("/problems", handlers.ListProblems)

	return r.Run()
}
