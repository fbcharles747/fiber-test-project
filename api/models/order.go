package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID    uint
	User      User
	ProductID uint `json:"product_id"`
	Product   Product
}
