package models

import (
	"time"

	"gorm.io/gorm"
)

type ListingStatus string

const (
	StatusActive   ListingStatus = "active"
	StatusSold     ListingStatus = "sold"
	StatusInactive ListingStatus = "inactive"
)

type Listing struct {
	gorm.Model  `json:"-"`
	ID          uint          `json:"id"          gorm:"primarykey"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Status      ListingStatus `json:"status"      gorm:"default:'active'"`
	Location    string        `json:"location"`
	ImageURL    string        `json:"image_url"`
	UserID      uint          `json:"user_id"`
	User        User          `json:"user,omitempty"     gorm:"foreignKey:UserID"`
	CategoryID  uint          `json:"category_id"`
	Category    Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time     `json:"created_at"`
}
