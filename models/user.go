package models

import "gorm.io/gorm"

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	gorm.Model   `json:"-"`
	ID           uint     `json:"id"    gorm:"primarykey"`
	Name         string   `json:"name"`
	Email        string   `json:"email" gorm:"unique"`
	Phone        string   `json:"phone"`
	Role         UserRole `json:"role"  gorm:"default:'user'"`
	PasswordHash string   `json:"-"`
}
