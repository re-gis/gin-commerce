package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/re-gis/gin-commerce/database"
	"gorm.io/gorm"
)

func CreateProduct(c *gin.Context) {
	// get the request
	var product database.Product
	var eProduct database.Product

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login first to access this feature"})
		return
	}

	var user database.User

	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	if user.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorised to perform the action..."})
		return
	}

	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if product.Description == "" || product.Name == "" || product.Price == 0 || product.StockQty == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all product details are required!"})
		return
	}

	if err := database.DB.Where("name =?", product.Name).First(&eProduct).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product already exists!"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error..."})
		return
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving the product!"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product saved successfully", "product": product})

}

func GetAllProducts(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login first to continue..."})
		return
	}

	var eUser database.User
	// get user from the database
	if err := database.DB.First(&eUser, userId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	if eUser.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised to perform this action..."})
		return
	}

	var products []database.Product
	if err := database.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Products not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Products fetched successfully...", "products": products})
}
