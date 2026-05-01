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

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		log.Printf("Incoming %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	r.POST("/api/register", handlers.CreateUser)
	r.POST("/api/login", handlers.Login)

	r.GET("/api/users/:id", handlers.GetUserByID)
	r.GET("/api/users/:id/listings", handlers.GetUserListings)

	r.GET("/api/categories", handlers.GetCategories)

	r.GET("/api/listings", handlers.GetListings)
	r.GET("/api/listings/:id", func(c *gin.Context) {
		id := c.Param("id")

		var listing models.Listing
		result := config.DB.Preload("User").Preload("Category").First(&listing, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
			return
		}

		reviews, err := GetReviewsByListingID(listing.ID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"listing": listing, "reviews": []interface{}{}})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"listing": listing,
			"reviews": reviews,
		})
	})

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/categories", handlers.CreateCategory)

		protected.POST("/listings", handlers.CreateListing)
		protected.PUT("/listings/:id", handlers.UpdateListing)
		protected.DELETE("/listings/:id", handlers.DeleteListing)
	}

	r.Run(":8080")
}
