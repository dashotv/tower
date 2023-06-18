package app

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/grimoire"
)

type DownloadFile struct {
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	MediumId primitive.ObjectID `json:"medium_id" bson:"medium_id"`
	Medium   *Medium            `json:"medium" bson:"medium"`
	Num      int                `json:"num" bson:"num"`
}

func NewDownloadFile() *DownloadFile {
	return &DownloadFile{}
}
