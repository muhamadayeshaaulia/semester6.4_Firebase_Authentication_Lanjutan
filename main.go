package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/muhamadayeshaaulia/gin-firebase-backend/config"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/routes"
)

func main() {
	//Load environment variables dari .env file
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variable sistem")
	}
	//inisialisasi firebase adminSDK
	config.InitFirebase()
	//inisialisasi database dan AutoMigrate
	config.InitDatabase()
	//setup gin router dengan semua router
	router := routes.SetupRouter()
	//port untuk menjalankan server
	port := os.Getenv(APP_PORT)
	if port == "" {
		port = "8080"
	}
	log.Printf("Server berjalan di http://localhost:%s", port)
	log.Printf("Health check: http://localhost:%s/v1/health", port)
}
