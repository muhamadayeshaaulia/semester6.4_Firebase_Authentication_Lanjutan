package config

import (
	"fmt"
	"os"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"

)

//DB instance GORM global yang di pakai di seluruh aplikasi
var DB *gorm.DB

func initDatabase(){
	//mengambil konfigurasi database dari environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//membuat format DSN (Data Source Name) untuk koneksi ke database
	//Format : user:pass@tcp(host:port)/dbname?params
	dsn := fmt.Sprint(
		"%s:%s@tcp(%s:%s)/%schartset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname,
	)
}