package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware memvalidasi backend JWT token di setiap request
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// MENGAMBIL TOKEN DARI HEADER AUTHORIZATION
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Authorization header tidak ditemukan",
				"error_code": "MISSING_TOKEN",
			})
			return
		}

		// Validasi format "Bearer <TOKEN>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Format token salah. Gunakan: Bearer <TOKEN>",
				"error_code": "INVALID_TOKEN_FORMAT",
			})
			return
		}

		tokenString := parts[1]

		// Parse dan verifikasi JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Memastikan algoritma yang dipakai adalah HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"message":    "Token tidak valid atau kadaluarsa",
				"error_code": "INVALID_TOKEN",
			})
			return
		}

		// Mengambil Claims dari Token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token Claims Tidak Valid",
			})
			return
		}
		
		var userID uint
		idFound := false

		// Cek beberapa kemungkinan nama key ID di dalam JWT (user_id, id, atau sub)
		keysToTry := []string{"user_id", "id", "sub"}
		
		for _, key := range keysToTry {
			if val, exists := claims[key]; exists && val != nil {
				// JWT biasanya menyimpan angka sebagai float64
				if floatVal, ok := val.(float64); ok {
					userID = uint(floatVal)
					idFound = true
					break
				}
				// Jika ID disimpan sebagai int (jarang di JWT tapi buat jaga-jaga)
				if intVal, ok := val.(int); ok {
					userID = uint(intVal)
					idFound = true
					break
				}
			}
		}

		if !idFound {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Data User ID tidak ditemukan dalam token. Silakan login ulang.",
			})
			return
		}

		// SET KE CONTEXT - Agar bisa diakses di handler manapun (termasuk AddToCart)
		c.Set("userID", userID)
		c.Set("email", claims["email"])
		c.Set("role", claims["role"])
		c.Set("firebase_uid", claims["firebase_uid"])

		c.Next() 
	}
}

// AdminOnly middleware untuk membatasi akses hanya untuk role admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success":    false,
				"message":    "Akses ditolak. Hanya Admin yang diizinkan!",
				"error_code": "FORBIDDEN",
			})
			return
		}
		c.Next()
	}
}