package utils

import (
	"backend-go/db"
	"context"
	"errors"
	"os"

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
func GetUserByToken(token string) (*User, error) {
    if token == "" || os.Getenv("JWT_SECRET") == "" {
        return nil, errors.New("invalid token or secret")
    }
    claims := &jwt.RegisteredClaims{}
    jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    if err != nil || !jwtToken.Valid {
        return nil, errors.New("invalid token")
    }
    userID, err := primitive.ObjectIDFromHex(claims.Subject)
    if err != nil {
        return nil, errors.New("invalid user id")
    }
    collection := db.GetCollection("users")
    var user User
    err = collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}


func GenerateToken(userID string, secret string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        Subject:   userID,
    })
    return token.SignedString([]byte(secret))
}