package controllers

import (
	"assignment-23/models/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	db *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{db: db}
}

func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	// Ambil JSON dari form field
	productJSON := ctx.PostForm("product")

	var product entity.Product

	// Unmarshal JSON ke struct Product
	// Jika JSON tidak valid, kembalikan error
	if err := json.Unmarshal([]byte(productJSON), &product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product JSON"})
		return
	}

	// Ambil dan validasi file
	// Gunakan ctx.FormFile untuk mendapatkan file dari form-data
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product image is required"})
		return
	}

	// Validasi ukuran file (maksimal 1MB)
	// Gunakan file.Size untuk mendapatkan ukuran file
	if file.Size > 1<<20 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File size must be under 1MB"})
		return
	}

	// Validasi ekstensi file (hanya JPG, JPEG, PNG)
	// Gunakan filepath.Ext untuk mendapatkan ekstensi file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG, JPEG, PNG files are allowed"})
		return
	}

	// Generate unique filename (optional tapi disarankan)
	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), product.Name)
	dst := filepath.Join("product-uploaded-images", filename+ext)

	// Simpan file duluan (di luar transaksi)
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Simpan ke DB dalam transaksi
	err = pc.db.Transaction(func(tx *gorm.DB) error {
		product.ProductImage = dst

		if err := tx.Create(&product).Error; err != nil {
			// Rollback manual ke file system
			_ = os.Remove(dst)
			return err
		}

		return nil
	})

	// Jika ada error, hapus file yang sudah disimpan
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product: " + err.Error()})
		return
	}

	// Jika berhasil, kembalikan response sukses
	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func (pc *ProductController) GetProducts(ctx *gin.Context) {
	var products []entity.Product

	if err := pc.db.Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (pc *ProductController) GetProductByID(ctx *gin.Context) {
	var product entity.Product
	id := ctx.Param("id")

	if err := pc.db.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (pc *ProductController) DownloadProductImageByProductName(c *gin.Context) {
	productname := c.Query("product_name")

	if productname == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Product name tidak boleh kosong!"})
		return
	}

	// Cari produk berdasarkan nama
	var product entity.Product
	if err := pc.db.Where("name = ?", productname).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	filePath := product.ProductImage
	filename := filepath.Base(filePath)

	fmt.Println("File Path:", filePath)
	fmt.Println("File Name:", filename)

	// Validasi apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Tambahkan Content-Disposition agar browser mengenali nama dan ekstensi file
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/octet-stream")

	c.File(filePath)
}

func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	var product entity.Product
	id := ctx.Param("id")

	// Cari produk berdasarkan ID
	if err := pc.db.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Invalid input": err.Error()})
		return
	}

	// Update produk dalam database
	if err := pc.db.Save(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func (pc *ProductController) DeleteProduct(ctx *gin.Context) {
	var product entity.Product
	id := ctx.Param("id")

	// Cari produk berdasarkan ID
	if err := pc.db.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Hapus produk dari database
	if err := pc.db.Delete(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
