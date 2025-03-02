package main

import (
	"dev-sinau/config"
	"dev-sinau/controllers"
	"dev-sinau/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("static", "./static")

	// Route untuk Register & Login
	r.GET("/register", controllers.RegisterPage)
	r.POST("/register", controllers.RegisterHandler(db))
	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.LoginHandler(db))
	r.POST("/logout", controllers.LogoutHandler())
	r.GET("/user", controllers.UpdateUserPage)

	// Proteksi halaman utama dengan middleware JWT
	protected := r.Group("")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/", controllers.TampilkanHalamanUtama(db))
		protected.POST("/buku", controllers.TambahBuku(db))
		protected.GET("/buku", controllers.LihatSemuaBuku(db))
		protected.PUT("/buku/:id", controllers.PerbaruiBuku(db))
		protected.DELETE("/buku/:id", controllers.HapusBuku(db))
		protected.PUT("/user", controllers.UpdateUser(db))

	}

	r.Run()
}
