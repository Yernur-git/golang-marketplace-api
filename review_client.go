package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

type ReviewResponse struct {
	ID        uint      `json:"id"`
	ListingID uint      `json:"listing_id"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

func GetReviewsByListingID(id uint) ([]ReviewResponse, error) {
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

	var reviews []ReviewResponse
	if err := json.Unmarshal(resp.Body(), &reviews); err != nil {
		return nil, err
	}

	return reviews, nil
}
