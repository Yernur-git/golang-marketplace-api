package handlers

import (
	"Marketplace-API/config"
	"Marketplace-API/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateListing(c *gin.Context) {
	var newListing models.Listing

	if err := c.ShouldBindJSON(&newListing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, _ := c.Get("user_id")
	newListing.UserID = userID.(uint)

	if newListing.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}
	if newListing.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price cannot be negative"})
		return
	}
	if newListing.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category_id is required"})
		return
	}

	if newListing.Status == "" {
		newListing.Status = models.StatusActive
	}

	result := config.DB.Create(&newListing)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create listing"})
		return
	}

	config.DB.Preload("User").Preload("Category").First(&newListing, newListing.ID)
	c.JSON(http.StatusCreated, newListing)
}

func GetListings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := config.DB.Model(&models.Listing{})

	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if location := c.Query("location"); location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var listings []models.Listing
	result := query.Preload("Category").
		Limit(limit).
		Offset(offset).
		Find(&listings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch listings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        listings,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": int(math.Ceil(float64(total) / float64(limit))),
	})
}

func GetListingByID(c *gin.Context) {
	id := c.Param("id")

	var listing models.Listing
	result := config.DB.Preload("User").Preload("Category").First(&listing, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	c.JSON(http.StatusOK, listing)
}

func UpdateListing(c *gin.Context) {
	id := c.Param("id")

	var listing models.Listing
	if err := config.DB.First(&listing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	var updateData struct {
		Title       string               `json:"title"`
		Description string               `json:"description"`
		Price       *float64             `json:"price"`
		Status      models.ListingStatus `json:"status"`
		Location    string               `json:"location"`
		CategoryID  *uint                `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if updateData.Title != "" {
		listing.Title = updateData.Title
	}
	if updateData.Description != "" {
		listing.Description = updateData.Description
	}
	if updateData.Price != nil {
		listing.Price = *updateData.Price
	}
	if updateData.Status != "" {
		listing.Status = updateData.Status
	}
	if updateData.Location != "" {
		listing.Location = updateData.Location
	}
	if updateData.CategoryID != nil {
		listing.CategoryID = *updateData.CategoryID
	}

	result := config.DB.Save(&listing)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update listing"})
		return
	}

	config.DB.Preload("User").Preload("Category").First(&listing, id)
	c.JSON(http.StatusOK, listing)
}

func DeleteListing(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&models.Listing{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing deleted successfully"})
}
