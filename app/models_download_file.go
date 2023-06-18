package app

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DownloadFile struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	MediumId primitive.ObjectID `json:"medium_id" bson:"medium_id"`
	Medium   *Medium            `json:"medium" bson:"medium"`
	Num      int                `json:"num" bson:"num"`
}

func NewDownloadFile() *DownloadFile {
	return &DownloadFile{}
}
