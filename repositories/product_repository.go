package repositories

import (
	"github.com/muhamadayeshaaulia/gin-firebase-backend/config"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"

)
type ProductRepository struct{}

func NewProductRepository() *ProductRepository{
	return &ProductRepository{}
}