package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/learn-qsharp/learn-qsharp-api/models"
	"github.com/lib/pq"
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

	tutorial := models.Tutorial{}

	sql := `
		SELECT id, title, description, difficulty, tags
		FROM tutorials
		WHERE tutorials.id = $1
	`

	err := db.QueryRow(ctx, sql, path.ID).Scan(&tutorial.ID, &tutorial.Title, &tutorial.Description, &tutorial.Difficulty,
		&tutorial.Tags)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	type tutorialLite struct {
		ID int `json:"id"`

		Title       string         `json:"title"`
		Description string         `json:"description"`
		Difficulty  string         `json:"difficulty"`
		Tags        pq.StringArray `json:"tags"`
	}

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

	tutorials := make([]tutorialLite, 0)
	for rows.Next() {
		tutorial := tutorialLite{}

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
