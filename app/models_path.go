package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Path struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Type      primitive.Symbol   `json:"type" bson:"type"`
	Remote    string             `json:"remote" bson:"remote"`
	Local     string             `json:"local" bson:"local"`
	Extension string             `json:"extension" bson:"extension"`
	Size      int                `json:"size" bson:"size"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
