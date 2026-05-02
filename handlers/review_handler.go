package handlers

import (
	"Marketplace-API/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func GetReviewsByListingID(c *gin.Context) {
	id := c.GetUint("listing_id")

	client := resty.New()

	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		log.Printf("[Resty] Requesting: %s %s", req.Method, req.URL)
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log.Printf("[Resty] Response Code: %d", resp.StatusCode())
		return nil
	})

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("http://review-service:8081/reviews?listing_id=%d", id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var reviews []models.ReviewResponse
	if err := json.Unmarshal(resp.Body(), &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
