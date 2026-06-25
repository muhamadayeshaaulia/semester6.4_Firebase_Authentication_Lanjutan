package models

import "gorm.io/gorm"

// Transaction mencatat tagihan/pesanan yang dibuat dari E-Commerce
type Transaction struct {
	gorm.Model
	UserID        uint    `gorm:"index" json:"user_id"`
	InvoiceID     string  `gorm:"uniqueIndex;size:100;not null" json:"invoice_id"`
	TotalAmount   float64 `gorm:"not null" json:"total_amount"`
	Status        string  `gorm:"size:20;default:'pending'" json:"status"` // pending, success, failed
	PaymentMethod string  `gorm:"size:50;default:'e-money'" json:"payment_method"`
}

// Request untuk membuat transaksi baru dari E-Commerce
type CreateTransactionRequest struct {
	TotalAmount float64 `json:"total_amount" binding:"required,gt=0"`
	ProductID   *uint   `json:"product_id,omitempty"`
	Quantity    *int    `json:"quantity,omitempty"`
}

// Request untuk aplikasi E-Money saat melakukan pembayaran
type PayTransactionRequest struct {
	InvoiceID string `json:"invoice_id" binding:"required"`
}
