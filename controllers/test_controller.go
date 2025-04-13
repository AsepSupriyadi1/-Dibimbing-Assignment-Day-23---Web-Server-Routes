package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type TestController struct{}

func NewTestController(db *gorm.DB) *TestController {
	return &TestController{}
}

func (ac *TestController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": "Server is runing on local port 8080"})
}
