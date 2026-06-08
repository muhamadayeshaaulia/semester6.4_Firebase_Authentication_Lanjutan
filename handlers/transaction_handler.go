package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/config"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"
)

func CreateTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoiceID := fmt.Sprintf("INV-%d", time.Now().Unix())

	tx := models.Transaction{
		UserID:      userID.(uint),
		InvoiceID:   invoiceID,
		TotalAmount: req.TotalAmount,
		Status:      "pending",
	}

	if err := config.DB.Create(&tx).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Setelah tagihan dibuat, kita juga mengosongkan keranjang (asumsi user mem-checkout semua di keranjang)
	config.DB.Where("user_id = ?", userID).Delete(&models.Cart{})

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction created",
		"data":    tx,
	})
}

func GetTransaction(c *gin.Context) {
	invoiceID := c.Param("invoice_id")

	var tx models.Transaction
	if err := config.DB.Where("invoice_id = ?", invoiceID).First(&tx).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tx,
	})
}

func GetUserTransactions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var transactions []models.Transaction
	config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}
