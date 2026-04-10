package main

import (
	"Marketplace-API/config"
	"Marketplace-API/handlers"
	"Marketplace-API/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	r.POST("/api/register", handlers.CreateUser)
	r.POST("/api/login", handlers.Login)

	r.GET("/api/users/:id", handlers.GetUserByID)
	r.GET("/api/users/:id/listings", handlers.GetUserListings)

	r.GET("/api/categories", handlers.GetCategories)

	r.GET("/api/listings", handlers.GetListings)
	r.GET("/api/listings/:id", handlers.GetListingByID)

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
