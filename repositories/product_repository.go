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

// findid mengambil satu produk berdasarkan id
func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
 var product models.Product
 result := config.DB.First(&product, id)
 return &product, result.Error
}
// Create menyimpan produk baru
func (r *ProductRepository) Create(product *models.Product) error {
 return config.DB.Create(product).Error
}
// Update memperbarui produk
func (r *ProductRepository) Update(product *models.Product) error {
 return config.DB.Save(product).Error
}
// Delete soft-delete produk (tidak hapus dari DB)
func (r *ProductRepository) Delete(id uint) error {
 return config.DB.Delete(&models.Product{}, id).Error
}