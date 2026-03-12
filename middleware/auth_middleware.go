package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"

)
// auth middleware memvalidasi backend jwt token di setiap request
//di pasang di route group yang memerlukan autentiksi
func AuthMiddleware() gin.HandlerFunc{
	return func (c *gin.Context){
		//MENGAMBIL TOKEN DARI HEADER AUTHORIZATION
		authHeader := c.GetHeader("Authorization")
		if authHeader == ""{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
				"success" : false,
				"message" : "Authorization header tidak ditemukan",
				"error_code": "MISSING_TOKEN",
			})
			return 
		}
	}
}