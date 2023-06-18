package app

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/grimoire"
)

type Download struct {
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	MediumId  primitive.ObjectID `json:"medium_id" bson:"medium_id"`
	Medium    Medium             `json:"medium" bson:"medium"`
	Auto      bool               `json:"auto" bson:"auto"`
	Multi     bool               `json:"multi" bson:"multi"`
	Force     bool               `json:"force" bson:"force"`
	Url       string             `json:"url" bson:"url"`
	ReleaseId string             `json:"release_id" bson:"tdo_id"`
	Thash     string             `json:"thash" bson:"thash"`
	Selected  string             `json:"selected" bson:"selected"`
	Status    string             `json:"status" bson:"status"`
	Files     []*DownloadFile    `json:"download_files" bson:"download_files"`
}

func NewDownload() *Download {
	return &Download{}
}
