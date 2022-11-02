package app

import (
	"time"

	"github.com/dashotv/grimoire"
)

type Feed struct {
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Name      string    `json:"name" bson:"name"`
	Url       string    `json:"url" bson:"url"`
	Source    string    `json:"source" bson:"source"`
	Type      string    `json:"type" bson:"type"`
	Active    bool      `json:"active" bson:"active"`
	Processed time.Time `json:"processed" bson:"processed"`
}

func NewFeed() *Feed {
	return &Feed{}
}
