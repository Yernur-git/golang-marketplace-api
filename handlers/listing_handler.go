package handlers

import (
	"Marketplace-API/config"
	"Marketplace-API/models"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

	var listing models.Listing
	if err := config.DB.First(&listing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	if listing.UserID != userID.(uint) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your listing"})
		return
	}

	config.DB.Delete(&listing)
	c.JSON(http.StatusOK, gin.H{"message": "Listing deleted successfully"})
}

func AdminDeleteListing(c *gin.Context) {
	id := c.Param("id")

	result := config.DB.Delete(&models.Listing{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing deleted by admin"})
}

func SearchListings(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var listings []models.Listing
	search := "%" + q + "%"
	result := config.DB.Preload("Category").
		Where("LOWER(title) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", search, search).
		Find(&listings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search listings"})
		return
	}

	c.JSON(http.StatusOK, listings)
}

func UpdateListingStatus(c *gin.Context) {
	id := c.Param("id")

	var listing models.Listing
	if err := config.DB.First(&listing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	var input struct {
		Status models.ListingStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	listing.Status = input.Status
	config.DB.Save(&listing)

	c.JSON(http.StatusOK, listing)
}

func GetRecentListings(c *gin.Context) {
	var listings []models.Listing
	result := config.DB.Preload("Category").
		Where("status = ?", models.StatusActive).
		Order("created_at DESC").
		Limit(8).
		Find(&listings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch listings"})
		return
	}

	c.JSON(http.StatusOK, listings)
}

func UploadListingImage(c *gin.Context) {
	id := c.Param("id")

	var listing models.Listing
	if err := config.DB.First(&listing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	userID, _ := c.Get("user_id")
	if listing.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your listing"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", listing.ID, time.Now().UnixNano(), ext)

	os.MkdirAll("uploads", 0755)
	dst := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	listing.ImageURL = "/uploads/" + filename
	config.DB.Save(&listing)

	c.JSON(http.StatusOK, gin.H{"image_url": listing.ImageURL})
}
