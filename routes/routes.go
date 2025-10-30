package routes

import (
	"go-crud/controllers"
	"go-crud/middleware"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	// ====== ROUTE PUBLIC ======
	e.GET("/", controllers.Home)
	e.POST("/register", controllers.Register)
	e.POST("/login", controllers.Login)

	// ====== ROUTE YANG BUTUH JWT ======
	api := e.Group("/api")
	api.Use(middleware.UseJWT())
	api.Use(middleware.AttachUser()) // <--- penting banget

	// ====== ROUTE PROFILE ======
	api.GET("/profile", controllers.Profile)

	// ====== ROUTE USERS ======
	users := api.Group("/users")
	users.GET("", controllers.GetAllUsers)       // hanya admin
	users.GET("/:id", controllers.GetUserByID)   // admin atau dirinya sendiri
	users.PUT("/:id", controllers.UpdateUser)    // admin atau dirinya sendiri
	users.DELETE("/:id", controllers.DeleteUser) // hanya admin

	// ====== ROUTE STORES ======
	stores := api.Group("/stores")
	stores.GET("", controllers.GetAllStores)       // hanya admin
	stores.GET("/my", controllers.GetMyStore)      // toko milik user login
	stores.GET("/:id", controllers.GetStoreByID)   // admin atau pemilik toko
	stores.PUT("/:id", controllers.UpdateStore)    // pemilik toko
	stores.DELETE("/:id", controllers.DeleteStore) // hanya admin

	// ====== ROUTE PRODUK (optional, bisa dipindah ke /api juga) ======
	e.GET("/products", controllers.ProductList)
	e.GET("/products/add", controllers.ProductForm)
	e.POST("/products/add", controllers.ProductAdd)
	e.GET("/products/edit/:id", controllers.ProductEditForm)
	e.POST("/products/edit/:id", controllers.ProductUpdate)
	e.GET("/products/delete/:id", controllers.ProductForm)
}
