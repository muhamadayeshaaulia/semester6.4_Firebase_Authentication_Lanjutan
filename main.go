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
}
