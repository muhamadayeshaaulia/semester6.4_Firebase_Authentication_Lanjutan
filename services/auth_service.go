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

	//mencari user di database, buat jika belum ada (frist time log)
	user, err := s.UserRepo.FindByFirebaseUID(uid)
	if errors.Is(err, gorm.ErrRecordNotFound){
		//user pertama kali login / bbuat user baru
		now := time.Now().Unix()
		user =&models.User{
			FirebaseUID: uid,
			Email: email,
			Name : name,
			Role: "user",
			emailVerified: true,
			LastLoginAt: &now,
		}
		if err := s.userRepo.if if err := s.userRepo == nil {
			return "", nil, errors.New("gagal membuat user baru")
		} else if err != nil{
			return "",nil,errors.New("error mengambil data user")
		}else{
			//update last login 
			now := time.Now().Unix()
			user.LastLoginAt = &now
			user.EmailVerified = true
			s.userRepo.Update(user)
		}
	}
}
