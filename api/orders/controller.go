package orders

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/re-gis/gin-commerce/database"
	"net/http"
)

type DeliverDetails struct {
	Order uint `json:"order"`
}

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
	// create the order
	var order database.Order
	//var orderItems []database.OrderItem

	order.Status = "PENDING"
	order.UserId = user.ID
	order.Cart = cart.ID
	if err := database.DB.Create(&order).Error; err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving the order"})
		return
	}

	if err := database.DB.Delete(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while clearing the cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully!"})
}

func Deliver(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login to continue"})
		return
	}

	// check if he is admin
	var user database.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	if user.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorised to perform this action!"})
		return
	}

	// bind the request
	var DeliverDetails DeliverDetails
	if err := c.BindJSON(&DeliverDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while binding the request body"})
		return
	}

	//get the order
	var order database.Order
	if err := database.DB.First(&order, DeliverDetails.Order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting the order"})
		return
	}

	//update the order
	order.Status = "DELIVERED"
	if err := database.DB.Save(&order).Error; err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating the order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order delivered successfully"})
}

func RejectOrder(c *gin.Context) {
	//get user
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login to continue"})
		return
	}

	var user database.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	if user.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorised to perform this action"})
		return
	}

	var DeleteDetails DeliverDetails
	if err := c.BindJSON(&DeleteDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while binding the data"})
		return
	}

	var order database.Order
	// delete the order
	if err := database.DB.Where("id = ?", DeleteDetails.Order).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found!"})
		return
	}

	if err := database.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting the order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order rejected successfully!"})

}
