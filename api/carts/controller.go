package carts

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/re-gis/gin-commerce/database"
)

type AddToCart struct {
	ProductId uint `json:"productId"`
	CartId    uint `json:"cartId"`
}

func AddItemToCart(c *gin.Context) {

	var addToCart AddToCart
	// get user
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login to continue..."})
		return
	}

	// bind the request
	if err := c.BindJSON(&addToCart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while binding the request body!"})
		return
	}

	var eCart database.Cart
	// check if cart existed
	if err := database.DB.First(&eCart, addToCart.CartId).Error; err != nil {
		// create the cart and save the cart items
		fmt.Print(userId)
	}
}

func RemoveItemFromCart(c *gin.Context) {

}
