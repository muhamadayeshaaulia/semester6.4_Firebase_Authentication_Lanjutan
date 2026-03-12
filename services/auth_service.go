package services

import (
	""
	"context"
	"errors"
	"github"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/muhamadayeshaaulia/gin-firebase-backend/config"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/repositories"
)

type AuthService struct {
	user *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{userRepo: repositories.NewUserRepository()}
}
