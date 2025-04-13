package dto

type UpdateStockByInventoryDTO struct {
	InventoryID uint `json:"inventory_id" binding:"required"`
	ProductID   uint `json:"product_id" binding:"required"`
	Qty         int  `json:"qty" binding:"required"` // positif untuk tambah, negatif untuk kurangi
}
