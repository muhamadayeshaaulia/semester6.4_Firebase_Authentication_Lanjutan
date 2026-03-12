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
