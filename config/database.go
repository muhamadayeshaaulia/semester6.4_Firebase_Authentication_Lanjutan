package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"
)

// DB instance GORM global yang di pakai di seluruh aplikasi
var DB *gorm.DB

func InitDatabase() *gorm.DB {
	//mengambil konfigurasi database dari environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//membuat format DSN (Data Source Name) untuk koneksi ke database
	//Format : user:pass@tcp(host:port)/dbname?params
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname,
	)

	//konfigurasi GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	//membuka koneksi ke database menggunakan GORM
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Gagal koneksi ke database : %v", err)
	}

	//setup connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan sql.DB : %v", err)
	}
	sqlDB.SetMaxOpenConns(25) // maksimal 25 koneksi yang terbuka
	sqlDB.SetMaxIdleConns(10) // maksimal 10 koneksi idle

	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Cart{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate gagal : %v", err)
	}
	log.Println("Database terhubung dan table sudah di migrate")
		return DB
}