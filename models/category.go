package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id"   gorm:"primarykey"`
	Name       string `json:"name" gorm:"unique"`
}
