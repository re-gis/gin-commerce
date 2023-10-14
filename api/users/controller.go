package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/re-gis/gin-commerce/database"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var newUser database.User

	// get the request
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if the user credentials are passed
	if newUser.Email == "" || newUser.Password == "" || newUser.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All credentials are requred!"})
		return
	}

	// hash the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
		return
	}

	newUser.Password = string(hashedPass)

	// save the user to database
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while registering user..."})
		return
	}

	// remove password from the user json to be sent to the frontend
	newUser.Password = ""
	c.JSON(http.StatusOK, newUser)
}
