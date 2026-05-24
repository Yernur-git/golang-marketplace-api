package models

import "time"

type ReviewResponse struct {
	ID        uint      `json:"id"`
	ListingID uint      `json:"listing_id"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
