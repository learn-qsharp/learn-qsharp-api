package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/learn-qsharp/learn-qsharp-api/models"
	"net/http"
)

func ShowProblem(c *gin.Context) {
	db := c.MustGet("db").(*pgxpool.Pool)
	ctx := c.MustGet("ctx").(context.Context)

	type pathStruct struct {
		ID *int `uri:"id" binding:"required"`
	}

	var path pathStruct
	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sql := `
		SELECT id, name, credits, body, template, difficulty, tags
		FROM problems
		WHERE problems.id = $1
	`

	problem := models.Problem{}

	err := db.QueryRow(ctx, sql, path.ID).Scan(&problem.ID, &problem.Name, &problem.Credits,
		&problem.Body, &problem.Template, &problem.Difficulty, &problem.Tags)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, problem)
}

func ListProblems(c *gin.Context) {
	db := c.MustGet("db").(*pgxpool.Pool)
	ctx := c.MustGet("ctx").(context.Context)

	sql := `
		SELECT id, name, difficulty, tags
		FROM problems
		ORDER BY id
	`

	rows, err := db.Query(ctx, sql)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	problems := make([]models.Problem, 0)
	for rows.Next() {
		problem := models.Problem{}

		err := rows.Scan(&problem.ID, &problem.Name, &problem.Difficulty, &problem.Tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		problems = append(problems, problem)
	}

	c.JSON(http.StatusOK, problems)
}
