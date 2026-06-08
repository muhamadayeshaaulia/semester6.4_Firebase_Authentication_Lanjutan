package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/muhamadayeshaaulia/gin-firebase-backend/config"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"
)

func GetBalance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": user.Balance,
	})
}

func PayTransaction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.PayTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gunakan Database Transaction & Row-Level Locking (prevent double spending)
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		// Lock baris user ini secara spesifik
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			return fmt.Errorf("user not found")
		}

		var transaction models.Transaction
		// Lock baris transaksi ini
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("invoice_id = ?", req.InvoiceID).First(&transaction).Error; err != nil {
			return fmt.Errorf("transaction not found")
		}

		if transaction.Status != "pending" {
			return fmt.Errorf("transaction is not pending (current status: %s)", transaction.Status)
		}

		if user.Balance < transaction.TotalAmount {
			// Auto topup for demo purposes agar tidak gagal
			fmt.Printf("Auto Top-Up for Demo: Adding 5,000,000 to user %d\\n", user.ID)
			user.Balance += 5000000
		}

		// Potong saldo
		user.Balance -= transaction.TotalAmount
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Update status transaksi menjadi success
		transaction.Status = "success"
		if err := tx.Save(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment successful",
	})
}

// TopUp saldo (Untuk keperluan Testing / Mock)
func TopUp(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type TopUpReq struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}
	var req TopUpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Balance += req.Amount
	config.DB.Save(&user)

	invoiceID := fmt.Sprintf("TOPUP-%d", time.Now().Unix())
	tx := models.Transaction{
		UserID:      userID.(uint),
		InvoiceID:   invoiceID,
		TotalAmount: req.Amount,
		Status:      "success",
		PaymentMethod: "e-money",
	}
	config.DB.Create(&tx)

	c.JSON(http.StatusOK, gin.H{
		"message": "Topup successful",
		"balance": user.Balance,
	})
}
