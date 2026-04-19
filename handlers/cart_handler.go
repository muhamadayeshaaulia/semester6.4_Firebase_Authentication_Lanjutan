package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"
	"gorm.io/gorm"
)

func AddToCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Struct lokal untuk menangkap request body
		// ProductID harus diawali Huruf Besar (Exported)
		var input struct {
			ProductID uint `json:"product_id" binding:"required"`
		}

		// Bongkar JSON dari Flutter
		if err := c.ShouldBindJSON(&input); err != nil {
			fmt.Printf("[DEBUG] Gagal Bind JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format data salah atau product_id kosong",
				"error":   err.Error(),
			})
			return
		}

		// Ambil UserID dari Middleware (Sudah dipastikan uint di middleware)
		val, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Sesi login habis"})
			return
		}
		userID := val.(uint)

		fmt.Printf("[DEBUG] BERHASIL! User: %d menambah Product: %d\n", userID, input.ProductID)
		var cart models.Cart
		
		// Cari apakah produk sudah ada di keranjang user ini
		err := db.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&cart).Error

		if err == gorm.ErrRecordNotFound {
			// Jika belum ada, buat record baru (CREATE)
			newCart := models.Cart{
				UserID:    userID,
				ProductID: input.ProductID,
				Quantity:  1,
			}
			if err := db.Create(&newCart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal simpan ke DB"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "Berhasil ditambah ke keranjang"})
		} else {
			// Jika sudah ada, tambahkan jumlahnya (UPDATE)
			cart.Quantity += 1
			db.Save(&cart)
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "Jumlah produk bertambah"})
		}
	}
}

// GetCart untuk mengambil daftar belanjaan user
func GetCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil userID dari token login
		val, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Silakan login ulang"})
			return
		}
		userID := val.(uint)

		var cartItems []models.Cart
		// Preload("Product") gunanya supaya relasi data produk ikut terambil (nama, harga, gambar)
		err := db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error

		if err != nil {
			fmt.Printf("[DEBUG] Gagal Get Cart: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data keranjang"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    cartItems,
		})
	}
}

func ReduceFromCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			ProductID uint `json:"product_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Input tidak valid"})
			return
		}

		val, _ := c.Get("userID")
		userID := val.(uint)

		var cart models.Cart
		// Cari item yang sesuai user dan produknya
		err := db.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&cart).Error

		if err == nil {
			if cart.Quantity > 1 {
				// Kurangi quantity jika lebih dari 1
				db.Model(&cart).Update("quantity", cart.Quantity-1)
				c.JSON(200, gin.H{"success": true, "message": "Jumlah dikurangi"})
			} else {
				// Hapus permanen jika sisa 1 lalu dikurangi
				db.Delete(&cart)
				c.JSON(200, gin.H{"success": true, "message": "Item dihapus dari keranjang"})
			}
		} else {
			c.JSON(404, gin.H{"error": "Item tidak ditemukan di keranjang"})
		}
	}
}