package models

import "gorm.io/gorm"

// User adalah model yang mapping ke tabel "users" di MySQL
// GORM otomatis plural nama struct -> nama tabel: User -> users

type User struct {
	gorm.Model           // Embed : ID, CreatedAt, DeletedAt(soft delete)
	FirebaseUID   string `gorm:"uniqueIndex;size:128;not null"json:"firebase_uid"`
	Email         string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Name          string `gorm:"size:100" json:"name"`
	Role          string `gorm:"size:20;default:user" json:"role"`
	EmailVerified bool   `gorm:"default:false"json:"email_verified"`
	LastLoginAt   *int64 `gorm:"index"json:"last_login_at,omitempty"`
}

// gorm.Model memberikan fields:
// ID uint (primary key, auto increment)
// CreatedAt time.Time (auto fill saat insert)
// UpdatedAt time.Time (auto fill saat update)
// DeletedAt gorm.DeletedAt (soft delete — row tidak benar-benar dihapus)
// Struct tag "gorm" mengontrol perilaku GORM:
// uniqueIndex = buat unique index di kolom ini
// size:128 = varchar(128)
// not null = kolom tidak boleh NULL
// default:user = nilai default kolom
// Struct tag "json" mengontrol serialisasi JSON:
// json:"firebase_uid" = nama key di JSON response
// json:"-" = field ini tidak dimasukkan ke JSON
// json:"...,omitempty" = skip jika nilnya zero/nil
