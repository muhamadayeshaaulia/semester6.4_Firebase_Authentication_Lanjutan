package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

)

// auth middleware memvalidasi backend jwt token di setiap request
// di pasang di route group yang memerlukan autentiksi
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//MENGAMBIL TOKEN DARI HEADER AUTHORIZATION
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Authorization header tidak ditemukan",
				"error_code": "MISSING_TOKEN",
			})
			return
		}
		//validasi format "Bearer <TOKEN>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"succes": false,
				"message": "Format token salah, Gunakan : Bearer <TOKEN>",
				"error_code":"INVALID_TOKE_FORMAT",
			})
			return
		}
		tokenString := parts[1]
		// parse dan verivikasi jwt token
		token, err := jwt.Parse(tokenString,func(token *jwt.Token)(interface{},error{
			//memastikan algo yang di pakai adalah hs256
			if _, ok := token.Method(*jwt.SigningMethodHMAC); !ok{
				return nil, jwt,jwt.ErrSignatureInvalid
			}
			return [] byte(os.Getenv("JWT_SECRET")), nil
		})
		
	)
	}
}
