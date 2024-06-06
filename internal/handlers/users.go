package handlers

import (
	"backend-go/db"
	"backend-go/internal/models"
	"backend-go/internal/utils"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// POST /api/v1/signup
func CreateUser(c *gin.Context) {
	// get user data from request
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	}
	collection := db.GetCollection("users")
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	result := collection.FindOne(ctx, bson.M{"email": user.Email})
	// check if user already exists
	if result.Err() == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}
	// add time field
	user.Date = time.Now()
	// insert user
    _, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error when adding user"})
	}
	// generate token
	token, err := utils.GenerateToken(user.ID.Hex(), os.Getenv("JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "token": token})
}


// POST /api/v1/login
func LoginUser(c *gin.Context) {
	var reqBody struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	collection := db.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	result := collection.FindOne(ctx, bson.M{"email": reqBody.Email})
	if result.Err() != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	err := result.Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}
	if reqBody.Password != user.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := utils.GenerateToken(user.ID.Hex(), os.Getenv("JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	}
	c.JSON(http.StatusOK, gin.H{"success": true,"message": "Login successful" ,"token": token})

}