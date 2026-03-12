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

// verifyFirebaseTokendari firebase
// memastikan email sudah terverivikasi, dan return backend jwt
func (s *AuthService) VerifyFirebaseToken(firebaseToken string) (string, *models.User, error) {
	//verivikasi firebase ID Token Ke server google
	token, err := config.FirebaseAuth.verifyIDToken(context.Background(), firebaseToken)
	if err != nil {
		return "", nil, errors.New("Firebase token tidak valid atau kadaluarsa")
	}
	//cek apakah emailsudah di verifikasi
	emailVerified, _ := token.Claims["email_verified"].(boll)
	if !emailVerified {
		return "", nil, errors.New("EMAIL_NOT_VERIFIED")
	}
	//mengambil data dari claims firebase token
	uid:= token.UID
	email, _ :=token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)

	
}
