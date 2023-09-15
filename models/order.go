package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model 
	UserID uint //Has One;One to Many
	BoughtItems []OrderItems //Has Many;Many to One
}

type OrderItems struct{
	gorm.Model
	OrderID uint //Has One;One to Many
	ProductID      uint
    BoughtProduct  Product `gorm:"foreignKey:ProductID"`
	Quantity int
	Address string
}

