package models

import "gorm.io/gorm"

type User struct {
	gorm.Model   `json:"-"`
	ID           uint   `json:"id"    gorm:"primarykey"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"unique"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"-"`
}
