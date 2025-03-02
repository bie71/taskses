package controllers

import (
	"dev-sinau/middleware"
	"dev-sinau/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Struktur request untuk register & login
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username    string `json:"username"`
	Password    string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

// Halaman Register
func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

// Halaman Login
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func UpdateUserPage(c *gin.Context) {
	c.HTML(http.StatusOK, "update-user.html", nil)
}

// Handler untuk register
func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
			return
		}

		// Cek apakah username sudah ada
		var existingUser models.User
		if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah digunakan"})
			return
		}

		// Hash password sebelum disimpan
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user := models.User{Username: req.Username, Password: string(hashedPassword)}
		db.Create(&user)

		c.JSON(http.StatusOK, gin.H{"message": "Registrasi berhasil"})
	}
}

// Handler untuk login
func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
			return
		}

		var user models.User
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
			return
		}

		// Generate token JWT
		tokenString, err := middleware.GenerateToken(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
			return
		}

		// Simpan token sebagai HTTP Cookie
		c.SetCookie("token", tokenString, 3600, "/", "", false, true)

		// Kirim response sukses
		c.JSON(http.StatusOK, gin.H{
			"message": "Login berhasil",
			"token":   tokenString,
		})
	}
}

func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
	}
}

// Update User
func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil username dari token JWT yang sudah divalidasi di middleware
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var user models.User
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}

		var req UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
			return
		}

		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password lama salah"})
			return
		}

		// Update username jika ada
		if req.Username != "" {
			user.Username = req.Username
		}

		// Update password jika ada
		if req.Password != "" {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			user.Password = string(hashedPassword)
		}

		// Simpan perubahan ke database
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User berhasil diperbarui"})
	}
}
