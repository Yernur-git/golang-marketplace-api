package handlers

import (
	"Marketplace-API/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func GetReviewsByListingID(id uint) ([]models.ReviewResponse, error) {
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
		return nil, err
	}

	var reviews []models.ReviewResponse
	if err := json.Unmarshal(resp.Body(), &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}
