package routes

import (
	"net/http"
	"text-adventure/constants"
	"text-adventure/models"
	"text-adventure/models/dbmodels"
	"text-adventure/utils"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	// In a real application, authenticate the user (this is just an example)
	body, err := utils.RequestBody(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userName := body["userName"].(string)
	password := body["password"].(string)

	// Check user credentials
	var user *dbmodels.User
	constants.DbClient.Where(&dbmodels.User{UserName: userName, Password: password}).First(&user)
	if user.Id == models.ZeroGuid() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.Id,
		"username": userName,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // Token expiration time
	})
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(constants.SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
