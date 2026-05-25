package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Review struct {
	ID        uint      `json:"id"`
	ListingID uint      `json:"listing_id"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	reviews   []Review
	mu        sync.Mutex
	idCounter uint = 1
)

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Header("Content-Type", "application/json")
		c.Next()
	})

	r.GET("/reviews", getReviews)
	r.POST("/reviews", createReview)

	r.Run(":8081")
}

func getReviews(c *gin.Context) {
	listingIDStr := c.Query("listing_id")
	if listingIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "listing_id is required"})
		return
	}

	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid listing_id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var result []Review
	for _, rev := range reviews {
		if rev.ListingID == uint(listingID) {
			result = append(result, rev)
		}
	}

	if result == nil {
		result = []Review{}
	}

	c.JSON(http.StatusOK, result)
}

func createReview(c *gin.Context) {
	var input struct {
		ListingID uint   `json:"listing_id" binding:"required"`
		Author    string `json:"author" binding:"required"`
		Rating    int    `json:"rating" binding:"required,min=1,max=5"`
		Comment   string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	rev := Review{
		ID:        idCounter,
		ListingID: input.ListingID,
		Author:    input.Author,
		Rating:    input.Rating,
		Comment:   input.Comment,
		CreatedAt: time.Now(),
	}
	idCounter++
	reviews = append(reviews, rev)

	c.JSON(http.StatusCreated, rev)
}
