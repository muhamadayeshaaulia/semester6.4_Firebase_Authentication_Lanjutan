package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/services"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{productService: services.NewProductService()}
}

// GET ALL PRODUCTS 
func (h *ProductHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")

	products, total, err := h.productService.GetAll(page, limit, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false, "message": "Gagal mengambil data produk",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    products,
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"limit":    limit,
			"per_page": limit,
		},
	})
}

// GET PRODUCT BY ID 
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}
	product, err := h.productService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Produk tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": product})
}

// CREATE PRODUCT (HYBRID IMAGE)
func (h *ProductHandler) Create(c *gin.Context) {
	var req models.CreateProductRequest

	// Bind menggunakan ShouldBind agar bisa terima JSON maupun Form-Data
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Cek upload file
	file, err := c.FormFile("image")
	if err == nil {
		// Skenario: Upload File
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + filepath.Base(file.Filename)
		savePath := "./public/uploads/" + filename

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal simpan file di server"})
			return
		}
		// Gunakan IP laptop Nafisah
		imageURL := "http://192.168.68.136:8080/uploads/" + filename
		req.ImageURL = imageURL
	} else if req.ImageURL == "" {
		// Skenario: Tidak ada file & tidak ada URL di body
		placeholderURL := "https://via.placeholder.com/400"
		req.ImageURL = placeholderURL
	}

	product, err := h.productService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat produk"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Produk berhasil dibuat", "data": product})
}

// UPDATE PRODUCT (WITH FILE CLEANUP)
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	// Ambil data lama untuk referensi hapus file
	oldProduct, err := h.productService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Produk tidak ditemukan"})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Cek jika ada gambar baru yang diupload
	file, err := c.FormFile("image")
	if err == nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "-" + filepath.Base(file.Filename)
		savePath := "./public/uploads/" + filename

		if err := c.SaveUploadedFile(file, savePath); err == nil {
			// Hapus file lama jika sebelumnya adalah file lokal
			if strings.Contains(oldProduct.ImageURL, "192.168.68.136") {
				oldFile := "./public/uploads/" + filepath.Base(oldProduct.ImageURL)
				os.Remove(oldFile)
			}
			imageURL := "http://192.168.68.136:8080/uploads/" + filename
			req.ImageURL = &imageURL
		}
	}

	product, err := h.productService.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal memperbarui produk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Produk diperbarui", "data": product})
}

// DELETE PRODUCT (WITH FILE CLEANUP)
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak valid"})
		return
	}

	// Ambil data untuk hapus file fisik
	product, err := h.productService.GetByID(uint(id))
	if err == nil && strings.Contains(product.ImageURL, "192.168.68.136") {
		fileName := filepath.Base(product.ImageURL)
		os.Remove("./public/uploads/" + fileName)
	}

	if err := h.productService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Produk tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Produk berhasil dihapus"})
}
