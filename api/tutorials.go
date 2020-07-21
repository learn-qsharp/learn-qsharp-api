package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/learn-qsharp/learn-qsharp-api/models"
	"net/http"
)

func ShowTutorial(c *gin.Context) {
	db := c.MustGet("db").(*pgx.Conn)
	ctx := c.MustGet("ctx").(context.Context)

	type pathStruct struct {
		ID int `uri:"id" binding:"required"`
	}

	var path pathStruct
	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sql := `
		SELECT id, title, credits, description, body, difficulty, tags
		FROM tutorials
		WHERE tutorials.id = $1
	`

	tutorial := models.Tutorial{}

	err := db.QueryRow(ctx, sql, path.ID).Scan(&tutorial.ID, &tutorial.Title, &tutorial.Credits, &tutorial.Description,
		&tutorial.Body, &tutorial.Difficulty, &tutorial.Tags)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, tutorial)
}

func ListTutorials(c *gin.Context) {
	db := c.MustGet("db").(*pgx.Conn)
	ctx := c.MustGet("ctx").(context.Context)

	sql := `
		SELECT id, title, description, difficulty, tags
		FROM tutorials
		ORDER BY id
	`

	rows, err := db.Query(ctx, sql)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tutorials := make([]models.Tutorial, 0)
	for rows.Next() {
		tutorial := models.Tutorial{}

		err := rows.Scan(&tutorial.ID, &tutorial.Title, &tutorial.Description, &tutorial.Difficulty,
			&tutorial.Tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		tutorials = append(tutorials, tutorial)
	}

	c.JSON(http.StatusOK, tutorials)
}
