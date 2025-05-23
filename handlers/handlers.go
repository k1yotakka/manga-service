package handlers

import (
	"github.com/gin-gonic/gin"
	"manga-catalog/client"
	"manga-catalog/database"
	"manga-catalog/models"
	"net/http"
	"strconv"
)

func GetMangaList(c *gin.Context) {
	var manga []models.Manga

	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	genre := c.Query("genre")

	limit, err1 := strconv.Atoi(limitStr)
	page, err2 := strconv.Atoi(pageStr)
	if err1 != nil || err2 != nil || limit <= 0 || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные параметры пагинации"})
		return
	}
	offset := (page - 1) * limit

	// Запрос с фильтрацией
	query := database.DB.Model(&models.Manga{})
	if genre != "" {
		query = query.Where("genre = ?", genre)
	}

	// Получаем общее количество
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при подсчёте total"})
		return
	}

	// Получаем текущую страницу
	if err := query.Limit(limit).Offset(offset).Find(&manga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	// Возвращаем данные + total
	c.JSON(http.StatusOK, gin.H{
		"data":  manga,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func CreateManga(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	genre := c.PostForm("genre")

	if title == "" || description == "" || genre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля обязательны"})
		return
	}

	manga := models.Manga{
		Title:       title,
		Description: description,
		Genre:       genre,
	}

	if err := database.DB.Create(&manga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении манги"})
		return
	}

	c.JSON(http.StatusCreated, manga)
}

func GetMangaByID(c *gin.Context) {
	id := c.Param("id")
	var manga models.Manga

	if err := database.DB.First(&manga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Манга не найдена"})
		return
	}

	c.JSON(http.StatusOK, manga)
}

func UpdateManga(c *gin.Context) {
	id := c.Param("id")
	var manga models.Manga

	if err := database.DB.First(&manga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Манга не найдена"})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	genre := c.PostForm("genre")

	if title != "" {
		manga.Title = title
	}
	if description != "" {
		manga.Description = description
	}
	if genre != "" {
		manga.Genre = genre
	}

	if err := database.DB.Save(&manga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении манги"})
		return
	}

	c.JSON(http.StatusOK, manga)
}

func DeleteManga(c *gin.Context) {
	id := c.Param("id")
	var manga models.Manga

	if err := database.DB.First(&manga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Манга не найдена"})
		return
	}

	if err := database.DB.Delete(&manga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении манги"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Манга успешно удалена"})
}

func AddToFavorites(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	mangaIDStr := c.Param("id")
	mangaID, err := strconv.Atoi(mangaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID манги"})
		return
	}

	var manga models.Manga
	if err := database.DB.First(&manga, mangaID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Манга не найдена"})
		return
	}

	var existing models.Favorite
	err = database.DB.Where("user_id = ? AND manga_id = ?", userID, mangaID).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Манга уже в избранном"})
		return
	}

	favorite := models.Favorite{
		UserID:  userID,
		MangaID: uint(mangaID),
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении в избранное"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Манга добавлена в избранное"})
}

func GetFavorites(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	user, err := client.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить данные пользователя"})
		return
	}

	var favorites []models.Favorite
	if err := database.DB.Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении избранного"})
		return
	}

	var mangaList []models.Manga
	for _, fav := range favorites {
		var manga models.Manga
		if err := database.DB.First(&manga, fav.MangaID).Error; err == nil {
			mangaList = append(mangaList, manga)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user":      user.Username,
		"favorites": mangaList,
	})
}

func RemoveFromFavorites(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	mangaIDStr := c.Param("id")
	mangaID, err := strconv.Atoi(mangaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID манги"})
		return
	}

	if err := database.DB.Where("user_id = ? AND manga_id = ?", userID, mangaID).Delete(&models.Favorite{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении из избранного"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Манга удалена из избранного"})
}
