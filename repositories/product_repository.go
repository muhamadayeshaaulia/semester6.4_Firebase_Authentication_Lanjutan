package repositories

import (
	"github.com/muhamadayeshaaulia/gin-firebase-backend/config"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"

)
type ProductRepository struct{}

func NewProductRepository() *ProductRepository{
	return &ProductRepository{}
}
//FindAll mengambil semua product aktif dengan pagination
func (r *ProductRepository) FindAll (page, limit int, category string)
([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := config.DB.Model(&models.Product{}). Where("is_actrive = ?", true)
	//filter by category jika ada
	if category != "" {
		query = query.Whare("category = ?", category)
	}
	//hitung total untuk pagination
	query.Count(&total)
	//ambil data dengan offset & limit
	offset := (page - 1) * limit
	result := query.Offset(offset).Limit(limit).Find(&products)
	return products, total, result.error
}