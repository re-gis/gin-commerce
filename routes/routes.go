package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/re-gis/gin-commerce/api/carts"
	"github.com/re-gis/gin-commerce/api/products"
	"github.com/re-gis/gin-commerce/api/users"
	"github.com/re-gis/gin-commerce/middleware"
)

func SetupRoutes(r *gin.Engine) *gin.Engine {
	r.POST("/users/login", users.LoginUser)
	r.POST("/users/register", users.RegisterUser)

	protected := r.Group("/")
	protected.Use(middleware.Authentication())
	{
		setupProductRoutes(protected)
		setupUserRoutes(protected)
		setupOrderRoutes(protected)
		setupCartRoutes(protected)
	}

	return r
}

func setupProductRoutes(rg *gin.RouterGroup) {
	productsRoute := rg.Group("/products")
	{
		productsRoute.POST("/create", products.CreateProduct)
		productsRoute.GET("/all", products.GetAllProducts)
		productsRoute.GET("/:id", products.GetOneProduct)
		productsRoute.DELETE("/delete/:id", products.DeleteProduct)
		productsRoute.PUT("/update/:id", products.UpdateProduct)

	}
}

func setupUserRoutes(rg *gin.RouterGroup) {
	usersRoute := rg.Group("/users")
	{
		usersRoute.GET("/all", users.GetAllUsers)
		usersRoute.PUT("/update/user/:id", users.UpdateUser)
		usersRoute.DELETE("/delete/myAccount", users.DeleteYouAccount)
		usersRoute.GET("/mine", users.GetYourAccount)
	}
}

func setupOrderRoutes(rg *gin.RouterGroup) {
	ordersRoute := rg.Group("/orders")
	{
		ordersRoute.GET("/order/all")
	}
}

func setupCartRoutes(rg *gin.RouterGroup) {
	cartsRoute := rg.Group("/carts")
	{
		cartsRoute.GET("/cart/all")
		cartsRoute.POST("/add", carts.AddItemToCart)
		cartsRoute.DELETE("/remove", carts.RemoveItemFromCart)
	}
}
