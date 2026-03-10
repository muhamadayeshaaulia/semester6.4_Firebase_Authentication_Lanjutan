package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name string `gorm:"size:200;not null;index" json:"name"`
	Description string `gorm:"type:text"json:"description"`
	Price float64 `gorm:"notr null"json:"price"`
	Stock int `gorm:"default:0"json:"stock"`
	Category string `grom:"size:100;index"json:"category"`
	ImageURL string `grom:"size:500"json:"image_url"`
	IsActive bool `gorm:"default:true;index"json:"is_active"`
}

// Request/Respons DTos (Data Transfer Objects)
// di pakai untuk validasi input HTTP request
type CreateProductRequest struct{
	Name string `json:"name"binding:"required,min=2,max=200"`
	Description string `json:"description"`
	Price float64 `json:"price" binding:"required,gt=0"`
	Stock int `json:"stock" binding:"min=0"`
	Category string `json:"category" binding:"required"`
	ImageURL string `json:"image_url"`
}

type UpdateProductRequest struct{
	Name *string `json:"name" binding:"omitempty, min=2"`
	Description *string `json:"description"`
	Price *float64 `json:"price" binding:"omitempty,gt=0"`
	Stock *int `json:"stock"binding:"imitempty,min=0"`
	Category *string `json:"category"`
	ImageURL *string `json:"image_url"`
}