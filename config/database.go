package config

import (
	"dev-sinau/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=sinau dbname=sk_dev port=9700 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// Cek apakah tabel sudah ada sebelum migrasi
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		log.Fatal("Gagal melakukan migrasi buku:", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Gagal melakukan migrasi user:", err)
	}

	return db
}
