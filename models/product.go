package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
	Price string `json:"price" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
	Category string `json:"category"`
	UserID    uint //One to Many
	OrderItems  []OrderItems  //Many to One
}