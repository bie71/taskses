package main

import (
	"dev-sinau/config"
	"dev-sinau/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", controllers.TampilkanHalamanUtama(db))
	r.POST("/buku", controllers.TambahBuku(db))
	r.GET("/buku", controllers.LihatSemuaBuku(db))
	r.PUT("/buku/:id", controllers.PerbaruiBuku(db))
	r.DELETE("/buku/:id", controllers.HapusBuku(db))

	r.Run(":8080")
}
