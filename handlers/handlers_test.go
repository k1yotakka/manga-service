package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"manga-catalog/database"
	"manga-catalog/handlers"
	"manga-catalog/middleware"
	"manga-catalog/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("DB_URL", "postgres://postgres:2705@localhost:5432/manga_test?sslmode=disable")
	database.ConnectDB()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())

	r.GET("/manga", handlers.GetMangaList)
	r.GET("/manga/:id", handlers.GetMangaByID)
	r.POST("/manga", handlers.CreateManga)
	r.PUT("/manga/:id", handlers.UpdateManga)
	r.DELETE("/manga/:id", handlers.DeleteManga)
	r.GET("/genres", handlers.GetAllGenres)
	r.GET("/genres/stats", handlers.GetGenresWithCount)
	r.POST("/manga/:id/comments", handlers.AddComment)
	r.GET("/manga/:id/comments", handlers.GetComments)
	r.POST("/manga/:id/favorite", handlers.AddToFavorites)
	r.DELETE("/manga/:id/favorite", handlers.RemoveFromFavorites)
	r.GET("/favorites", handlers.GetFavorites)

	return r
}

func generateToken(userID uint, role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(middleware.JwtSecret)
	return tokenString
}

func TestGetMangaListSuccess(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/manga", nil)
	req.Header.Set("Authorization", "Bearer "+generateToken(1, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetMangaListInvalidPagination(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/manga?limit=-1", nil)
	req.Header.Set("Authorization", "Bearer "+generateToken(1, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateMangaMissingFields(t *testing.T) {
	r := setupRouter()
	body := `{"title": "", "description": "", "genre": "", "cover": ""}`
	req, _ := http.NewRequest("POST", "/manga", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateToken(1, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateMangaSuccess(t *testing.T) {
	r := setupRouter()
	body := `{"title": "Naruto", "description": "Ninja story", "genre": "Action", "cover": "cover.jpg"}`
	req, _ := http.NewRequest("POST", "/manga", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateToken(1, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestGetMangaByIDNotFound(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/manga/999999", nil)
	req.Header.Set("Authorization", "Bearer "+generateToken(1, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestUpdateMangaNotFound(t *testing.T) {
	r := setupRouter()
	body := `{"title": "NewTitle", "description": "NewDesc", "genre": "Drama", "cover": "cover.jpg"}`
	req, _ := http.NewRequest("PUT", "/manga/999999", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateToken(1, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestDeleteMangaNotFound(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("DELETE", "/manga/999999", nil)
	req.Header.Set("Authorization", "Bearer "+generateToken(1, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestGetAllGenres(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/genres", nil)
	req.Header.Set("Authorization", "Bearer "+generateToken(1, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetGenresWithCount(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/genres/stats", nil)
	req.Header.Set("Authorization", "Bearer "+generateToken(1, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestAddCommentInvalidJSON(t *testing.T) {
	r := setupRouter()
	manga := models.Manga{Title: "M1", Description: "D", Genre: "G", Cover: "C.jpg"}
	database.DB.Create(&manga)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/manga/%d/comments", manga.ID), bytes.NewBuffer([]byte(`invalid`)))
	req.Header.Set("Authorization", "Bearer "+generateToken(2, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestAddCommentSuccess(t *testing.T) {
	r := setupRouter()
	manga := models.Manga{Title: "M2", Description: "D", Genre: "G", Cover: "C.jpg"}
	database.DB.Create(&manga)

	comment := map[string]string{"text": "Good one!"}
	body, _ := json.Marshal(comment)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/manga/%d/comments", manga.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateToken(2, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestGetComments(t *testing.T) {
	r := setupRouter()
	manga := models.Manga{Title: "M3", Description: "D", Genre: "G", Cover: "C.jpg"}
	database.DB.Create(&manga)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/manga/%d/comments", manga.ID), nil)
	req.Header.Set("Authorization", "Bearer "+generateToken(2, "user"))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestAddAndRemoveFromFavorites(t *testing.T) {
	r := setupRouter()
	manga := models.Manga{Title: "FavM", Description: "Desc", Genre: "Genre", Cover: "C.jpg"}
	database.DB.Create(&manga)
	token := generateToken(3, "user")

	// Add to favorites
	req, _ := http.NewRequest("POST", fmt.Sprintf("/manga/%d/favorite", manga.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Get favorites
	req, _ = http.NewRequest("GET", "/favorites", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Remove from favorites
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/manga/%d/favorite", manga.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
