package dto

type CreateOrderDTO struct {
	OrderItems []CreateOrderItemDTO `json:"order_items" binding:"required,dive"`
}

type CreateOrderItemDTO struct {
	ProductID uint `json:"product_id" binding:"required"`
	Qty       int  `json:"qty" binding:"required,gt=0"`
}

type OrderItemResponse struct {
	ProductID uint `json:"product_id"`
	Qty       int  `json:"qty"`
}
