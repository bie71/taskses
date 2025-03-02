package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Secret key untuk signing token (Gantilah dengan yang lebih aman)
var secretKey = []byte("supersecretkey")

// GenerateToken untuk autentikasi user

// Claims untuk JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Generate Token JWT
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour) // Token berlaku 1 jam
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Buat token dengan metode HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Validasi Token JWT
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	// Cek apakah token sudah expired
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token sudah expired")
	}

	return claims, nil
}

// JWTAuth Middleware untuk proteksi endpoint
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// Cek Header Authorization
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Jika tidak ada di header, cek cookie
			var err error
			token, err = c.Cookie("token")
			if err != nil || token == "" {
				if c.GetHeader("Accept") == "application/json" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				} else {
					c.Redirect(http.StatusFound, "/login")
				}
				c.Abort()
				return
			}
		}

		// Validasi token
		claims, err := ValidateToken(token)
		if err != nil {
			if c.GetHeader("Accept") == "application/json" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			} else {
				c.Redirect(http.StatusFound, "/login")
			}
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
