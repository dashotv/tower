package app

import (
	"time"

	"github.com/dashotv/grimoire"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Watch struct {
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Username  string             `json:"username" bson:"username"`
	Player    string             `json:"player" bson:"player"`
	WatchedAt time.Time          `json:"watched_at" bson:"watched_at"`
	MediumId  primitive.ObjectID `json:"medium_id" bson:"medium_id"`
	Medium    *Medium            `json:"medium" bson:"medium"`
}

func NewWatch() *Watch {
	return &Watch{}
}
