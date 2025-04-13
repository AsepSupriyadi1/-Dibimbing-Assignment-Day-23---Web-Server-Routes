package controllers

import (
	"assignment-23/models/dto"
	"assignment-23/models/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryController struct {
	db *gorm.DB
}

func NewInventoryController(db *gorm.DB) *InventoryController {
	return &InventoryController{db: db}
}

func (ic *InventoryController) CreateInventory(ctx *gin.Context) {
	var inventory entity.Inventory

	// Validasi masukan dari pengguna
	if err := ctx.ShouldBindJSON(&inventory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid input!",
		})
		return
	}

	if err := ic.db.Create(&inventory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to create inventory",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Inventory created successfully"})
}
func (ic *InventoryController) GetInventories(ctx *gin.Context) {
	var inventories []entity.Inventory

	if err := ic.db.Find(&inventories).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventories"})
		return
	}

	ctx.JSON(http.StatusOK, inventories)
}

func (ic *InventoryController) GetProductStock(ctx *gin.Context) {
	productID := ctx.Param("id")

	var totalStock int64
	if err := ic.db.Model(&entity.Inventory{}).
		Where("product_id = ?", productID).
		Select("SUM(stock)").Scan(&totalStock).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "failed to fetch stock"})
		return
	}

	ctx.JSON(200, gin.H{
		"product_id":  productID,
		"total_stock": totalStock,
	})
}

func (ic *InventoryController) UpdateStockByInventory(ctx *gin.Context) {
	var input dto.UpdateStockByInventoryDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var inventory entity.Inventory
	if err := ic.db.First(&inventory, input.InventoryID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
		return
	}

	if inventory.ProductID != input.ProductID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "product_id does not match with inventory"})
		return
	}

	newStock := inventory.Stock + input.Qty
	if newStock < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "stock cannot be negative"})
		return
	}

	inventory.Stock = newStock
	if err := ic.db.Save(&inventory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update stock"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "stock updated successfully",
		"inventory_id":  inventory.ID,
		"product_id":    inventory.ProductID,
		"current_stock": inventory.Stock,
	})
}
