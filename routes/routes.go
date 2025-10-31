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
	api.Use(middleware.AttachUser()) 

	// ====== ROUTE PROFILE ======
	api.GET("/profile", controllers.Profile)

	// ====== ROUTE USERS ======
	users := api.Group("/users")
	{
		users.GET("", controllers.GetAllUsers)       
		users.GET("/:id", controllers.GetUserByID)   
		users.PUT("/:id", controllers.UpdateUser)    
		users.DELETE("/:id", controllers.DeleteUser) 
	}

	// ====== ROUTE STORES ======
	toko := api.Group("/toko")
	{
		toko.GET("", controllers.GetAllToko)       
		toko.GET("/my", controllers.GetMyToko)      
		toko.GET("/:id", controllers.GetTokoByID)   
		toko.PUT("/:id", controllers.UpdateToko)    
		toko.DELETE("/:id", controllers.DeleteToko) 
	}

	// ====== ROUTE PRODUCTS ======
	// products := api.Group("/products")
	// {
	// 	products.GET("", controllers.GetAllProducts)     
	// 	products.GET("/:id", controllers.GetProductByID) 
	// 	products.POST("", controllers.CreateProduct)     
	// 	products.PUT("/:id", controllers.UpdateProduct)  
	// 	products.DELETE("/:id", controllers.DeleteProduct) 
	// }

	// ====== ROUTE ALAMAT ======
	alamat := api.Group("/alamat")
	{
		alamat.GET("/my", controllers.GetMyAlamat)   
		alamat.POST("", controllers.CreateAlamat)      
		alamat.GET("/:id", controllers.GetAlamatByID)      
		alamat.PUT("/:id", controllers.UpdateAlamat)   
		alamat.DELETE("/:id", controllers.DeleteAlamat) 
	}

	// ====== Province & City routes =====
	provcity := e.Group("/provcity")
	{
		provcity.GET("/listprovinces", controllers.GetListProvinces)
		provcity.GET("/listcities/:province_id", controllers.GetListCities)
		provcity.GET("/detailprovince/:id", controllers.GetDetailProvince)
		provcity.GET("/detailcity/:id", controllers.GetDetailCity)
	}

	// ====== ROUTE KATEGORI (admin only) ======
	// categories := api.Group("/categories")
	// {
	// 	categories.GET("", controllers.GetAllCategories)      
	// 	categories.POST("", controllers.CreateCategory)       
	// 	categories.PUT("/:id", controllers.UpdateCategory)    
	// 	categories.DELETE("/:id", controllers.DeleteCategory) 
	// }

	// ====== ROUTE TRANSAKSI ======
	// transactions := api.Group("/transactions")
	// {
	// 	transactions.GET("/my", controllers.GetMyTransactions)       
	// 	transactions.POST("", controllers.CreateTransaction)         
	// 	transactions.GET("/:id", controllers.GetTransactionByID)     
	// 	transactions.PUT("/:id/status", controllers.UpdateStatusTransaction) 
	// }
}
