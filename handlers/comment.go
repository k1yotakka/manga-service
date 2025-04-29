package handlers

import (
	"manga-catalog/database"
	"manga-catalog/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	mangaIDStr := c.Param("id")
	mangaID, err := strconv.Atoi(mangaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID манги"})
		return
	}

	var body struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Текст обязателен"})
		return
	}

	comment := models.Comment{
		MangaID:   uint(mangaID),
		UserID:    userID,
		Text:      body.Text,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func GetComments(c *gin.Context) {
	mangaID := c.Param("id")

	var comments []models.Comment
	err := database.DB.Where("manga_id = ?", mangaID).Order("created_at DESC").Find(&comments).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении комментариев"})
		return
	}

	c.JSON(http.StatusOK, comments)
}
