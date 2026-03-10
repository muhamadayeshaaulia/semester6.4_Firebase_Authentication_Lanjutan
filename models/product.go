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
