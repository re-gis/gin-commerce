package routes

import (
	"github.com/gin-gonic/gin"
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
	products := rg.Group("/products")
	{
		products.GET("/product/all")
	}
}

func setupUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("/user/all")
	}
}

func setupOrderRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders")
	{
		orders.GET("/order/all")
	}
}

func setupCartRoutes(rg *gin.RouterGroup) {
	carts := rg.Group("/carts")
	{
		carts.GET("/cart/all")
	}
}
