package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/learn-qsharp/learn-qsharp-api/models"
	"net/http"
)

func ShowTutorial(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	type pathStruct struct {
		ID int `uri:"id" binding:"required"`
	}

	var path pathStruct
	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	tutorial := models.Tutorial{}

	if err := db.First(&tutorial, path.ID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, tutorial)
}

func ListTutorials(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	tutorials := make([]models.Tutorial, 0)
	if err := db.Order("id").Find(&tutorials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tutorials)
}
