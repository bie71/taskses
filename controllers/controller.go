package controllers

import (
	"dev-sinau/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TampilkanHalamanUtama(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var daftarBuku []models.Book
		db.Find(&daftarBuku)
		c.HTML(http.StatusOK, "index.html", gin.H{"buku": daftarBuku})
	}
}

func TambahBuku(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buku models.Book
		if err := c.ShouldBindJSON(&buku); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&buku)
		c.JSON(http.StatusCreated, buku)
	}
}

func LihatSemuaBuku(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var daftarBuku []models.Book
		db.Find(&daftarBuku)
		c.JSON(http.StatusOK, daftarBuku)
	}
}

func PerbaruiBuku(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var buku models.Book
		if err := db.First(&buku, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Buku tidak ditemukan"})
			return
		}
		if err := c.ShouldBindJSON(&buku); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&buku)
		c.JSON(http.StatusOK, buku)
	}
}

func HapusBuku(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var buku models.Book
		if err := db.First(&buku, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Buku tidak ditemukan"})
			return
		}
		db.Delete(&buku)
		c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dihapus"})
	}
}
