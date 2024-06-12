package utils

import (
	"backend-go/db"
	"backend-go/internal/models"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User structure (adjust according to your actual User model)
type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    CartData map[string]int     `bson:"cartData,omitempty"`
}

// Function to get user by token
func GetUserByToken(token string) (*models.User, error) {
	if token == "" || os.Getenv("JWT_SECRET") == "" {
		return nil, errors.New("invalid token or secret")
	}
	claims := &jwt.RegisteredClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		fmt.Println("Token parsing error:", err)
		return nil, errors.New("invalid token")
	}
	if !jwtToken.Valid {
		fmt.Println("Token is not valid")
		return nil, errors.New("invalid token")
	}
	userID, err := primitive.ObjectIDFromHex(claims.Subject)
	if err != nil {
		fmt.Println("Invalid user ID in token claims:", err)
		return nil, errors.New("invalid user id")
	}
	collection := db.GetCollection("users")
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		fmt.Println("User lookup error:", err)
		return nil, err
	}
	return &user, nil
}



func GenerateToken(userID string, secret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}