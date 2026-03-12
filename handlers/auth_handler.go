package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/muhamadayeshaaulia/gin-firebase-backend/services"

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
	//Verifikasi via service
	jwtToken, user, err := h.authService.VerifyFirebaseToken(req.FirebaseToken)
	if err != nil {
		// Bedakan error email belum verify vs error lainnya
		if err.Error() == "EMAIL_NOT_VERIFIED" {
			c.JSON(http.StatusForbidden, gin.H{
				"success":    false,
				"message":    "Email belum diverifikasi. Cek inbox email Anda.",
				"error_code": "EMAIL_NOT_VERIFIED",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    err.Error(),
				"error_code": "INVALID_FIREBASE_TOKEN",
			})
		}
		return
	}
	//Return Backend JWT + data user
	expireHours := 24
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"data": gin.H{
			"access_token": jwtToken,
			"token_type":   "Bearer",
			"expires_in":   expireHours * 3600,
			"user": gin.H{
				"id":             user.ID,
				"firebase_uid":   user.FirebaseUID,
				"email":          user.Email,
				"name":           user.Name,
				"role":           user.Role,
				"email_verified": user.EmailVerified,
				"created_at":     user.CreatedAt.Format(time.RFC3339),
			},
		},
	})
}
