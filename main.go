package main

import (
	"Marketplace-API/config"
	"Marketplace-API/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	r.POST("/api/users", handlers.CreateUser)
	r.GET("/api/users/:id", handlers.GetUserByID)
	r.GET("/api/users/:id/listings", handlers.GetUserListings)

	r.POST("/api/categories", handlers.CreateCategory)
	r.GET("/api/categories", handlers.GetCategories)

	r.POST("/api/listings", handlers.CreateListing)
	r.GET("/api/listings", handlers.GetListings)
	r.GET("/api/listings/:id", handlers.GetListingByID)
	r.PUT("/api/listings/:id", handlers.UpdateListing)
	r.DELETE("/api/listings/:id", handlers.DeleteListing)

	r.Run(":8080")
}
