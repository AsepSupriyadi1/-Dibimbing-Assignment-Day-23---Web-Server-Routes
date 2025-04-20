package dto

// Request dari user untuk mendownload sebuah file
type DownloadProductRequest struct {
	ProductName string `json:"product_name" binding:"required"`
}
