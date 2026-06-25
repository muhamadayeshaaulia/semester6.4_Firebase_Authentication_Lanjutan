package models

import "gorm.io/gorm"

// Transaction mencatat tagihan/pesanan yang dibuat dari E-Commerce
type Transaction struct {
	gorm.Model
	UserID        uint    `gorm:"index" json:"user_id"`
	InvoiceID     string  `gorm:"uniqueIndex;size:100;not null" json:"invoice_id"`
	TotalAmount   float64 `gorm:"not null" json:"total_amount"`
	Status        string  `gorm:"size:20;default:'pending'" json:"status"` // pending, success, failed
	PaymentMethod string            `gorm:"size:50;default:'e-money'" json:"payment_method"`
	Items         []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
}

// TransactionItem mencatat detail produk yang dibeli
type TransactionItem struct {
	gorm.Model
	TransactionID uint    `gorm:"index" json:"transaction_id"`
	ProductID     uint    `json:"product_id"`
	Product       Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"` // Harga satuan saat dibeli
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
