package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name        string             `bson:"name" json:"name"`
    Image       string             `bson:"image" json:"image"`
    Category    string             `bson:"category" json:"category"`
    NewPrice    float64            `bson:"new_price" json:"new_price"`
    OldPrice    float64            `bson:"old_price" json:"old_price"`
    Description string             `bson:"description" json:"description"`
    Date        time.Time          `bson:"date" json:"date"`
    Available   bool               `bson:"available" json:"available" default:"true"`
}
