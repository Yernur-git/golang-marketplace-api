package main

import (
	"Marketplace-API/config"
	"Marketplace-API/handlers"
	"Marketplace-API/middleware"
	"Marketplace-API/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "http://localhost:5173" || origin == "http://localhost:5174" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Use(func(c *gin.Context) {
		log.Printf("Incoming %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	r.POST("/api/register", handlers.CreateUser)
	r.POST("/api/login", handlers.Login)

	r.GET("/api/users/:id", handlers.GetUserByID)
	r.GET("/api/users/:id/listings", handlers.GetUserListings)

	r.GET("/api/categories", handlers.GetCategories)
	r.GET("/api/categories/:id", handlers.GetCategoryByID)
	r.GET("/api/categories/:id/listings", handlers.GetCategoryListings)

	r.GET("/api/listings", handlers.GetListings)
	r.GET("/api/listings/recent", handlers.GetRecentListings)
	r.GET("/api/listings/search", handlers.SearchListings)
	r.GET("/api/listings/:id", func(c *gin.Context) {
		id := c.Param("id")

		var listing models.Listing
		result := config.DB.Preload("User").Preload("Category").First(&listing, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
			return
		}

		reviews, err := handlers.GetReviewsByListingID(listing.ID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"listing": listing, "reviews": []interface{}{}})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"listing": listing,
			"reviews": reviews,
		})
	})

	r.Static("/uploads", "./uploads")

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", handlers.GetMe)
		protected.PUT("/users/profile", handlers.UpdateProfile)

		protected.POST("/categories", handlers.CreateCategory)

		protected.POST("/listings", handlers.CreateListing)
		protected.PUT("/listings/:id", handlers.UpdateListing)
		protected.PATCH("/listings/:id/status", handlers.UpdateListingStatus)
		protected.DELETE("/listings/:id", handlers.DeleteListing)
		protected.POST("/listings/:id/image", handlers.UploadListingImage)
	}

	return r
}

func main() {
	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.Listing{})

	r := SetupRouter()
	r.Run(":8080")
}
