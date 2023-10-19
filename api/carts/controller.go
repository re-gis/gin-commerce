package carts

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/re-gis/gin-commerce/database"
	"gorm.io/gorm"
)

type AddToCart struct {
	ProductId uint `json:"productId"`
	CartId    uint `json:"cartId"`
	Quantity  int  `json:"quantity"`
}

type RemoveItemFromCartDtls struct {
	ProductId uint `json:"productId"`
	Quantity  int  `json:"quantity"`
}

func AddItemToCart(c *gin.Context) {
	var addToCart AddToCart

	// get user
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login to continue..."})
		return
	}

	// Ensure the userId can be converted to uint
	uid, ok := userId.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// bind the request
	if err := c.BindJSON(&addToCart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while binding the request body!"})
		return
	}

	if addToCart.CartId == 0 || addToCart.ProductId == 0 || addToCart.Quantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All details are required!"})
		return
	}

	// check if product is available
	var product database.Product
	if err := database.DB.First(&product, addToCart.ProductId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found!"})
		return
	}

	if product.StockQty < addToCart.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product in stock is not enough"})
		return
	}

	var eCart database.Cart
	// check if cart existed
	if err := database.DB.First(&eCart, addToCart.CartId).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			// create the cart and save the cart items
			var cart database.Cart
			cart.UserId = uid
			// Save the cart
			if err := database.DB.Create(&cart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a cart!"})
				return
			}
		}
	}

	// Assuming you'll also need to create a cart item at this point
	var existingCartItem database.CartItem
	var cartItem database.CartItem

	// if cartItem exists increment the quantity
	if err := database.DB.Where("cart_id = ? and product_id =?", eCart.ID, addToCart.ProductId).First(&existingCartItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			cartItem.CartId = eCart.ID
			cartItem.ProductId = addToCart.ProductId
			cartItem.Quantity = addToCart.Quantity // Or another value based on your logic

			if err := database.DB.Create(&cartItem).Error; err != nil {
				fmt.Print(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart!"})
				return
			}

			var user database.User
			if err := database.DB.First(&user, userId).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
				return
			}

			if user.Cart == 0 {
				user.Name = user.Name
				user.Password = user.Password
				user.Email = user.Email
				user.Cart = eCart.ID
				// saving user cart
				if err := database.DB.Save(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating the user's cart"})
					return
				}
			}
		} else {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking the cart item"})
			return
		}

	} else {
		existingCartItem.Quantity += addToCart.Quantity // or any logic to adjust the quantity
		if err := database.DB.Save(&existingCartItem).Error; err != nil {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item quantity!"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added successfully!"})
}

func RemoveItemFromCart(c *gin.Context) {
	// get the user cart
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

	var cart database.Cart
	if err := database.DB.First(&cart, user.Cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	var removeItemFromCartDtls RemoveItemFromCartDtls
	// bind the request body
	if err := c.BindJSON(&removeItemFromCartDtls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while binding the request body"})
		return
	}

	/* logic to remove the cart item
	and if it was with quantity== 1 deletes
	it from cart else decrement
	by the quantity to remove */

	var cartItem database.CartItem
	if err := database.DB.Where("cart_id = ? AND product_id = ?", cart.ID, removeItemFromCartDtls.ProductId).First(&cartItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving cart item"})
		}
		return
	}

	if cartItem.Quantity == 1 || cartItem.Quantity == removeItemFromCartDtls.Quantity {
		// Delete the cart item
		if err := database.DB.Delete(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart!"})
			return
		}
	} else if cartItem.Quantity > removeItemFromCartDtls.Quantity {
		// Reduce the cart item's quantity
		cartItem.Quantity -= removeItemFromCartDtls.Quantity
		if err := database.DB.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item's quantity!"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove more items than exist in cart!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated successfully!"})

}
