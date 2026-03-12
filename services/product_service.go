package services

import (
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/repositories"

)

type ProductService struct {
	ProductRepo *repositories.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{ProductRepo: repositories.NewProductRepository()}
}
func (s *ProductService) GetAll(page, limit int, category string) ([]models.Product, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	return s.ProductRepo.FindAll(page, limit, category)
}
func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	return s.ProductRepo.FindByID(id)
}
func (s *ProductService) Create(req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
	}
	err := s.ProductRepo.Create(product)
	return product, err
}
func (s *ProductService) Update(id uint, req *models.UpdateProductRequest) (*models.Product, error) {
	product, err := s.ProductRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	//update hanya field yang di kirim (pointer nil tidak di update)
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.ImageURL != nil {
		product.ImageURL = *req.ImageURL
	}

	err = s.ProductRepo.Update(product)
	return product,err
}

