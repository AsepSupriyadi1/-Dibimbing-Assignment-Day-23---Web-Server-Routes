package routes

import (
	"assignment-23/controllers"
	"assignment-23/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	testController := controllers.NewTestController(db)
	authController := controllers.NewAuthController(db)
	productController := controllers.NewProductController(db)
	inventoryController := controllers.NewInventoryController(db)
	orderController := controllers.NewOrderController(db)

	api := router.Group("/api")
	{

		api.GET("/test", testController.Ping)

		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
		}

		protected := api.Group("/protected")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/test", testController.Ping)

			// Product routes
			product := protected.Group("/product")
			{
				product.GET("/", productController.GetProducts)
				product.GET("/:id", productController.GetProductByID)
				product.POST("/", productController.CreateProduct)
				product.PUT("/:id", productController.UpdateProduct)
				product.DELETE("/:id", productController.DeleteProduct)
			}

			// Inventory routes
			inventory := protected.Group("/inventory")
			{
				inventory.GET("/", inventoryController.GetInventories)
				inventory.POST("/", inventoryController.CreateInventory)
				inventory.GET("/stock/:id", inventoryController.GetProductStock)
				inventory.POST("/update-stock", inventoryController.UpdateStockByInventory)
			}

			// Order routes
			order := protected.Group("/order")
			{
				order.GET("/:id", orderController.GetOrderByID)
				order.POST("/", orderController.CreateOrder)
			}
		}
	}
}
