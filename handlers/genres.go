package handlers

import (
	"github.com/gin-gonic/gin"
	"manga-catalog/database"
	"manga-catalog/models"
	"net/http"
)

func GetAllGenres(c *gin.Context) {
	var genres []string
	err := database.DB.
		Model(&models.Manga{}).
		Distinct("genre").
		Pluck("genre", &genres).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить жанры"})
		return
	}

	c.JSON(http.StatusOK, genres)
}

func GetGenresWithCount(c *gin.Context) {
	type Result struct {
		Genre string
		Count int
	}
	var result []Result

	err := database.DB.
		Model(&models.Manga{}).
		Select("genre, COUNT(*) as count").
		Group("genre").
		Scan(&result).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при подсчёте жанров"})
		return
	}

	response := make(map[string]int)
	for _, r := range result {
		response[r.Genre] = r.Count
	}

	c.JSON(http.StatusOK, response)
}
