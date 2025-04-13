package controllers

import (
	"assignment-23/models/dto"
	"assignment-23/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

func (ac *AuthController) Login(ctx *gin.Context) {

	var loginRequest dto.LoginRequest

	// Validasi masukan dari pengguna
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi kredensial
	if loginRequest.Email != os.Getenv("DB_ACCOUNT_EMAIL") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	// Periksa password
	if loginRequest.Password != os.Getenv("DB_ACCOUNT_PASSWORD") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(1)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
