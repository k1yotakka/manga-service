package main

import (
	"manga-catalog/database"
	"manga-catalog/handlers"
	"manga-catalog/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.AuthMiddleware())

	api := r.Group("/api")
	{
		api.GET("/manga", handlers.GetMangaList)
		api.POST("/manga", handlers.CreateManga)
		api.GET("/manga/:id", handlers.GetMangaByID)
		api.PUT("/manga/:id", handlers.UpdateManga)
		api.DELETE("/manga/:id", handlers.DeleteManga)

		api.GET("/genres", handlers.GetAllGenres)
		api.GET("/genres/stats", handlers.GetGenresWithCount)

		api.POST("/manga/:id/comments", handlers.AddComment)
		api.GET("/manga/:id/comments", handlers.GetComments)

		api.POST("/manga/:id/favorite", handlers.AddToFavorites)
		api.GET("/favorites", handlers.GetFavorites)
		api.DELETE("/manga/:id/favorite", handlers.RemoveFromFavorites)
	}

	r.Run(":8080")
}
