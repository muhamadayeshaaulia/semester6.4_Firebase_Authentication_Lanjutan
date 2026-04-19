package models

import "gorm.io/gorm"

type Product struct {
    gorm.Model
    Name        string  `gorm:"size:200;not null;index" json:"name"`
    Description string  `gorm:"type:text" json:"description"`
    Price       float64 `gorm:"not null" json:"price"` 
    Stock       int     `gorm:"default:0" json:"stock"`
    Category    string  `gorm:"size:100;index" json:"category"`
    ImageURL    string  `gorm:"size:500" json:"image_url"`
    IsActive    bool    `gorm:"default:true;index" json:"is_active"`
}

type CreateProductRequest struct {
    Name        string  `json:"name" form:"name" binding:"required,min=2,max=200"`
    Description string  `json:"description" form:"description"`
    Price       float64 `json:"price" form:"price" binding:"required,gt=0"`
    Stock       int     `json:"stock" form:"stock" binding:"min=0"`
    Category    string  `json:"category" form:"category" binding:"required"`
    ImageURL    string  `json:"image_url" form:"image_url"`
}

type UpdateProductRequest struct {
    Name        *string  `json:"name" binding:"omitempty,min=2"` 
    Description *string  `json:"description"`
    Price       *float64 `json:"price" binding:"omitempty,gt=0"`
    Stock       *int     `json:"stock" binding:"omitempty,min=0"` 
    Category    *string  `json:"category"`
    ImageURL    *string  `json:"image_url"`
}