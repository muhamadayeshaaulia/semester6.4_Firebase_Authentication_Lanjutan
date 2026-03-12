package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/muhamadayeshaaulia/gin-firebase-backend/services"
	"net/http"
	"time"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{authService: services.NewAuthService()}
}

// POST /auth/verify-token
// Terima Firebase ID Token → verifikasi → return Backend JWT
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	// 1. Parse request body
	var req struct {
		FirebaseToken string `json:"firebase_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "firebase_token wajib diisi",
		})
		return
	}
}
