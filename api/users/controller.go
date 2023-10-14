package users

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/re-gis/gin-commerce/database"
	"github.com/re-gis/gin-commerce/utils"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = os.Getenv("JWT_SECRET")

func RegisterUser(c *gin.Context) {
	var newUser database.User
	var eUser database.User

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

	// check if user aleady exists
	rs := database.DB.Where("email =?", newUser.Email).First(&eUser)
	if rs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists!"})
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

func LoginUser(c *gin.Context) {
	var cred utils.Credentials
	var user database.User

	// get the request
	if err := c.BindJSON(&cred); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while binding the request"})
		return
	}

	// check if all credentials are provided
	if cred.Email == "" || cred.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All credentials are required!"})
		return
	}

	// get user from database
	if err := database.DB.Where("email = ?", cred.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password!"})
		return
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password!"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
