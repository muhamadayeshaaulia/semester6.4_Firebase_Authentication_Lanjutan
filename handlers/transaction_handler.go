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

	// Update Stock and Sold
	var transactionItems []models.TransactionItem

	if req.ProductID != nil && req.Quantity != nil && *req.Quantity > 0 {
		// Beli Langsung (Hanya 1 produk, tidak mengosongkan keranjang)
		var product models.Product
		if err := config.DB.First(&product, *req.ProductID).Error; err == nil {
			transactionItems = append(transactionItems, models.TransactionItem{
				TransactionID: tx.ID,
				ProductID:     product.ID,
				Quantity:      *req.Quantity,
				Price:         product.Price,
			})
			product.Stock -= *req.Quantity
			if product.Stock < 0 {
				product.Stock = 0
			}
			product.Sold += *req.Quantity
			config.DB.Save(&product)
		}
	} else {
		// Beli dari keranjang
		var cartItems []models.Cart
		config.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems)
		for _, item := range cartItems {
			var product models.Product
			if err := config.DB.First(&product, item.ProductID).Error; err == nil {
				transactionItems = append(transactionItems, models.TransactionItem{
					TransactionID: tx.ID,
					ProductID:     product.ID,
					Quantity:      item.Quantity,
					Price:         product.Price,
				})
				product.Stock -= item.Quantity
				if product.Stock < 0 {
					product.Stock = 0
				}
				product.Sold += item.Quantity
				config.DB.Save(&product)
			}
		}
		// Setelah tagihan dibuat, kita juga mengosongkan keranjang
		config.DB.Where("user_id = ?", userID).Delete(&models.Cart{})
	}

	if len(transactionItems) > 0 {
		config.DB.Create(&transactionItems)
	}

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
	config.DB.Preload("Items").Preload("Items.Product").Where("user_id = ?", userID).Order("created_at desc").Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

func UpdateTransactionStatus(c *gin.Context) {
	invoiceID := c.Param("invoice_id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tx models.Transaction
	if err := config.DB.Where("invoice_id = ?", invoiceID).First(&tx).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	tx.Status = req.Status
	if err := config.DB.Save(&tx).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated", "data": tx})
}
