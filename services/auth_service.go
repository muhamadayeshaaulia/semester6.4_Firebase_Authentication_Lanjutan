package services

import (
	"context"
	"errors"
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
	userRepo *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{userRepo: repositories.NewUserRepository()}
}

// verifyFirebaseTokendari firebase
func (s *AuthService) VerifyFirebaseToken(firebaseToken string) (string, *models.User, error) {
    token, err := config.FirebaseAuth.VerifyIDToken(context.Background(), firebaseToken)
    if err != nil {
        return "", nil, errors.New("Firebase token tidak valid atau kadaluarsa")
    }

    emailVerified, _ := token.Claims["email_verified"].(bool)
    if !emailVerified {
        return "", nil, errors.New("EMAIL_NOT_VERIFIED")
    }

    uid := token.UID
    email, _ := token.Claims["email"].(string)
    name, _ := token.Claims["name"].(string)

    user, err := s.userRepo.FindByFirebaseUID(uid)
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        now := time.Now().Unix()
        user = &models.User{
            FirebaseUID:   uid,
            Email:         email,
            Name:          name,
            Role:          "user",
            EmailVerified: true,
            LastLoginAt:   &now,
        }
        if err := s.userRepo.Create(user); err != nil {
            return "", nil, errors.New("gagal membuat user baru")
        }
    } else if err != nil {
        return "", nil, errors.New("error mengambil data user")
    } else {
        // Kalau user sudah ada, kita harus update status verifikasinya dan last login
        now := time.Now().Unix()
        user.LastLoginAt = &now
        user.EmailVerified = true // SINKRONISASI STATUS DI SINI
        
        // Simpan perubahan ke database EliteBook
        if err := s.userRepo.Update(user); err != nil {
            return "", nil, errors.New("gagal update data user")
        }
    }

    // Generate JWT (Sekarang claims email_verified di JWT pasti true)
    jwtToken, err := s.generateJWT(user)
    if err != nil {
        return "", nil, errors.New("gagal membuat token")
    }
    return jwtToken, user, nil
}
func (s *AuthService) CreateUserInMySQL(uid, email, name string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("firebase_uid = ?", uid).First(&user).Error

	if err == nil {
		return &user, nil 
	}
	newUser := models.User{
		FirebaseUID: uid,
		Email:       email,
		Name:        name,
		Role:        "user",
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

// generate token jwt dengan payload user
func (s *AuthService) generateJWT(user *models.User) (string, error) {
	expireHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if expireHours == 0 {
		expireHours = 24
	}
	//payload yang di simpan dalam token
	claims := jwt.MapClaims{
		"sub":            user.ID,
		"firebase_uid":   user.FirebaseUID,
		"email":          user.Email,
		"name":           user.Name,
		"role":           user.Role,
		"email_verified": user.EmailVerified,
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(time.Hour * time.Duration(expireHours)).Unix(),
	}
	// membuat token dengan algo HS256 dan secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
