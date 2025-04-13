package entity

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CustomerName string      `json:"customer_name"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Qty       int
	Order     Order   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product   Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
