package orders

import (
	"github.com/gin-gonic/gin"
	"github.com/re-gis/gin-commerce/database"
	"net/http"
)

func PlaceOrder(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login to continue"})
		return
	}

	// get the user cart
	var user database.User
	var cart database.Cart

	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := database.DB.First(&cart, user.Cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found!"})
		return
	}

	// get all cart items in the cart
	var cartItems []database.CartItem
	if err := database.DB.Where("cart_id =?", cart.ID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No cart items found!"})
		return
	}
	if len(cartItems) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No cart items found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cartitems": cartItems})
}
