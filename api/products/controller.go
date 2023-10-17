package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/re-gis/gin-commerce/database"
	"gorm.io/gorm"
)

type ProductUpdate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	StockQty    int     `json:"stock_qty"`
}

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
	var products []database.Product
	if err := database.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Products not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Products fetched successfully...", "products": products})
}

func GetOneProduct(c *gin.Context) {
	// get the id
	productId := c.Param("id")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No id provided!"})
		return
	}

	// get the product
	var product database.Product

	if err := database.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product fetched successfully", "product": product})
}

func DeleteProduct(c *gin.Context) {
	productId := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login to continue!"})
		return
	}

	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product id not provided!"})
		return
	}

	var user database.User
	// check the role
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	if user.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorised to perform this action..."})
		return
	}

	var product database.Product
	// get product
	if err := database.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found!"})
		return
	}

	// delete the product
	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting the product!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully..."})
}

func UpdateProduct(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login first to continue"})
		return
	}

	// get product id
	productId := c.Param("id")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product id not provided!"})
		return
	}

	var user database.User
	// get user role
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	if user.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorised to perform this action..."})
		return
	}

	var product database.Product
	// get product
	if err := database.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found!"})
		return
	}

	var productUpdateDetails ProductUpdate
	// get request body
	if err := c.ShouldBindJSON(&productUpdateDetails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error..."})
		return
	}

	

	// update the product
	if productUpdateDetails.Description != "" {
		product.Description = productUpdateDetails.Description
	}
	if productUpdateDetails.Name != "" {
		product.Name = productUpdateDetails.Name
	}
	if productUpdateDetails.Price != 0 {
		product.Price = productUpdateDetails.Price
	}

	if productUpdateDetails.StockQty != 0 {
		product.StockQty = productUpdateDetails.StockQty
	}

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating the product!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": product})
}
