package middlewares

import (
	"net/http"
	"text-adventure/models/dbmodels"
	"text-adventure/utils"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const secretKey = "my_secret_key"

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
	zeroUser := &dbmodels.User{}
	utils.DbClient.Where(&dbmodels.User{UserName: userName, Password: password}).First(&user)
	if user.Id == zeroUser.Id {
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
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Stop further processing if unauthorized
			return
		}

		// Set the token claims to the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next() // Proceed to the next handler if authorized
	}
}
