package entity

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	ProductID uint    `json:"product_id"`
	Product   Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Stock     int     `json:"stock"`
	Location  string  `json:"location"`
}
