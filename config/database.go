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