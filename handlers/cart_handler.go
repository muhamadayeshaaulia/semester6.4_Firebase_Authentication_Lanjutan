package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"
)

func AddToCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			ProductID uint `json:"product_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ambil UserID dari Middleware Auth (yang sudah kita buat sebelumnya)
		userID := c.MustGet("userID").(uint)

		var cart models.Cart
		// Cek apakah produk ini sudah ada di keranjang user tersebut
		err := db.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&cart).Error

		if err == gorm.ErrRecordNotFound {
			// Jika belum ada, buat baru
			newCart := models.Cart{
				UserID:    userID,
				ProductID: input.ProductID,
				Quantity:  1,
			}
			db.Create(&newCart)
			c.JSON(http.StatusOK, gin.H{"message": "Berhasil ditambah ke keranjang"})
		} else {
			// Jika sudah ada, update jumlahnya
			cart.Quantity += 1
			db.Save(&cart)
			c.JSON(http.StatusOK, gin.H{"message": "Jumlah produk bertambah di keranjang"})
		}
	}
}
