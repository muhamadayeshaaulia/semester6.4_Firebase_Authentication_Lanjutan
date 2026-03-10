package repositories

import (
	"github.com/muhamadayeshaaulia/gin-firebase-backend/config"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/models"

)
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindByFirebaseUID mencari berdasarkan firebase UID
func (r *UserRepository) FindByFirebaseUID(uid string) (*models.User, error){
	var user models.User
	result := config.DB.Where("Firebase_Uid = ?", uid).First(&user)
	if result.Error != nil{
		return nil, result.Error
	}
	return &user,nil
}
//find mail mencari user berdasarkan email yang terdaftar di firebase
func (r *UserRepository) FindByEmail(email string)(*models.User, error){
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	return &user,result.Error
}

//create menyimpan user baru ke database
