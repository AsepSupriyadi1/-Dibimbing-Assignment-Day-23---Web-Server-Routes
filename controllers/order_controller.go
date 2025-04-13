package controllers

import (
	"assignment-23/models/dto"
	"assignment-23/models/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	db *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{db: db}
}

func (oc *OrderController) CreateOrder(ctx *gin.Context) {
	var input dto.CreateOrderDTO

	// Validate input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := oc.db.Begin()

	// Create order
	var order entity.Order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	var orderItems []entity.OrderItem

	// Loop through order items
	for _, item := range input.OrderItems {
		var inventory entity.Inventory

		// Check if product exists in inventory
		if err := tx.Where("product_id = ?", item.ProductID).First(&inventory).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %d not found in inventory", item.ProductID)})
			return
		}

		if inventory.Stock < item.Qty {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("not enough stock for product %d", item.ProductID)})
			return
		}

		// Update stock
		inventory.Stock -= item.Qty
		if err := tx.Save(&inventory).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update inventory"})
			return
		}

		// Add item to orderItems
		orderItems = append(orderItems, entity.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Qty:       item.Qty,
		})
	}

	// Create order items after order has an ID
	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order items"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	// Prepare response
	var responseItems []dto.OrderItemResponse
	for _, item := range orderItems {
		responseItems = append(responseItems, dto.OrderItemResponse{
			ProductID: item.ProductID,
			Qty:       item.Qty,
		})
	}

	// Send response
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "order created & stock updated successfully",
		"data": gin.H{
			"order_id":    order.ID,
			"order_items": responseItems,
		},
	})
}

// Get Order Detail Including Order Items By ID
func (oc *OrderController) GetOrderByID(ctx *gin.Context) {
	var order entity.Order
	id := ctx.Param("id")

	if err := oc.db.Preload("OrderItems").First(&order, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	ctx.JSON(http.StatusOK, order)
}
