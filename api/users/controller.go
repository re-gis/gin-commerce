package users

import (
	"net/http"
	"os"

	// "strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/re-gis/gin-commerce/database"
	"github.com/re-gis/gin-commerce/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtKey = os.Getenv("JWT_SECRET")

type UserUpdate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

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
	if err := database.DB.Where("email = ?", newUser.Email).First(&eUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists!"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
		return
	}

	// hash the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
		return
	}

	newUser.Password = string(hashedPass)

	// save role if present
	if newUser.Role != "" {
		newUser.Role = "admin"
	}

	// save the user to database
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while registering user..."})
		return
	}

	// make a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    newUser.ID,
		"email": newUser.Email,
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error..."})
		return
	}

	// remove password from the user json to be sent to the frontend
	newUser.Password = ""
	c.JSON(http.StatusOK, gin.H{"user": newUser, "token": tokenString})
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

// update the user
func UpdateUser(c *gin.Context) {
	var updateUserDetails UserUpdate

	// bind the request
	if err := c.ShouldBindJSON(&updateUserDetails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error..."})
		return
	}

	// get user id
	userId := c.Param("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user id provided"})
		return
	}

	// get the user from the database
	var user database.User
	if err := database.DB.Where("id =?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user!"})
		return
	}

	// update the user
	if updateUserDetails.Email != "" {
		user.Email = updateUserDetails.Email
	}

	if updateUserDetails.Password != "" {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(updateUserDetails.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password!"})
			return
		}
		user.Password = string(hashedPass)
	}

	if updateUserDetails.Name != "" {
		user.Name = updateUserDetails.Name
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving the user..."})
		return
	}

	// return the user
	user.Password = ""
	c.JSON(http.StatusOK, user)

}

func DeleteYouAccount(c *gin.Context) {
	// get user id
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login first to access this feature!"})
		return
	}

	var user database.User
	// get user from database
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// delete user
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while deleting the user account..."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User account deleted successfully!"})

}

func GetAllUsers(c *gin.Context) {
	// get the user id from token to check user role
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login first to access the following feature..."})
		return
	}

	var user database.User
	// get user
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	// get the role
	if user.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorised to access this feature..."})
		return
	}

	// get the users
	var users []database.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Users not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "users fetched!", "users": users})
}

func GetYourAccount(c *gin.Context) {
	// get id
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login to access your account profile..."})
		return
	}

	var user database.User
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	c.JSON(http.StatusOK, user)
}
