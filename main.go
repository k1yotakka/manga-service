package main

import (
	"github.com/gin-contrib/cors"
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
	r.Static("/uploads", "./uploads")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/manga", handlers.GetMangaList)
		api.GET("/manga/:id", handlers.GetMangaByID)
		api.GET("/genres", handlers.GetAllGenres)
		api.GET("/genres/stats", handlers.GetGenresWithCount)
		api.GET("/manga/:id/comments", handlers.GetComments)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/manga", handlers.CreateManga)
		protected.PUT("/manga/:id", handlers.UpdateManga)
		protected.DELETE("/manga/:id", handlers.DeleteManga)
		protected.POST("/manga/:id/comments", handlers.AddComment)
		protected.POST("/manga/:id/favorite", handlers.AddToFavorites)
		protected.GET("/favorites", handlers.GetFavorites)
		protected.DELETE("/manga/:id/favorite", handlers.RemoveFromFavorites)
	}

	r.Run(":8080")
}
