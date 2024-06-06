package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name        string             `bson:"name" json:"name"`
	Email		string			   `bson:"email" json:"email"`
	Password	string			   `bson:"password" json:"password"`
	CartData	map[string]int	   `bson:"cartData" json:"cartData"`
    Date        time.Time          `bson:"date" json:"date"`
}