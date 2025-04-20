package entity

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name         string      `json:"name" gorm:"not null"  binding:"required"`
	Description  string      `json:"description"`
	Price        float64     `json:"price" gorm:"not null"  binding:"required"`
	Category     string      `json:"category"  binding:"required"`
	ProductImage string      `json:"product_image"`
	Inventories  []Inventory `gorm:"foreignKey:product_id"`
	OrderItems   []OrderItem `gorm:"foreignKey:product_id"`
}

//type ProductCategory struct {
//	gorm.Model
//	Name    string    `json:"name" gorm:"not null"`
//	Product []Product `gorm:"foreignKey:product_category_id"`
//}
