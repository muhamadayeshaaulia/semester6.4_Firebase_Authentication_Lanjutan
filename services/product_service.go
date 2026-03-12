package services

import (
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/repositories"

)

type ProductService struct {
	ProductRepo *repositories.ProductRepository
}
func NewProductService() *	ProductService{
	return &ProductService{ProductRepo: repositories.NewProductRepository()}
}
func (s *ProductService) GetAll(page, limit int, category string)([]models.Product, int64,error){
	if page <= 0 { page  = 1}
	if limit <= 0 || limit> 100 {limit =10}
}
