package app

import (
	"time"

	"github.com/dashotv/grimoire"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Episode struct {
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Type           string             `json:"type" bson:"_type"`
	SeriesId       primitive.ObjectID `json:"series_id" bson:"series_id"`
	Kind           primitive.Symbol   `json:"kind" bson:"kind"`
	Source         string             `json:"source" bson:"source"`
	SourceId       string             `json:"source_id" bson:"source_id"`
	Title          string             `json:"title" bson:"title"`
	Description    string             `json:"description" bson:"description"`
	Slug           string             `json:"slug" bson:"slug"`
	Text           []string           `json:"text" bson:"text"`
	Display        string             `json:"display" bson:"display"`
	Directory      string             `json:"directory" bson:"directory"`
	Search         string             `json:"search" bson:"search"`
	SeasonNumber   int                `json:"season_number" bson:"season_number"`
	EpisodeNumber  int                `json:"episode_number" bson:"episode_number"`
	AbsoluteNumber int                `json:"absolute_number" bson:"absolute_number"`
	SearchParams   struct {
		Type       string `json:"type" bson:"type"`
		Verified   bool   `json:"verified" bson:"verified"`
		Group      string `json:"group" bson:"group"`
		Author     string `json:"author" bson:"author"`
		Resolution int    `json:"resolution" bson:"resolution"`
		Source     string `json:"source" bson:"source"`
		Uncensored bool   `json:"uncensored" bson:"uncensored"`
		Bluray     bool   `json:"bluray" bson:"bluray"`
	} `json:"search_params" bson:"search_params"`
	Active      bool      `json:"active" bson:"active"`
	Downloaded  bool      `json:"downloaded" bson:"downloaded"`
	Completed   bool      `json:"completed" bson:"completed"`
	Skipped     bool      `json:"skipped" bson:"skipped"`
	Watched     bool      `json:"watched" bson:"watched"`
	Broken      bool      `json:"broken" bson:"broken"`
	Unwatched   int       `json:"unwatched" bson:"unwatched"`
	ReleaseDate time.Time `json:"release_date" bson:"release_date"`
	Paths       []Path    `json:"paths" bson:"paths"`
	Cover       string    `json:"cover" bson:"cover"`
	Background  string    `json:"background" bson:"background"`
}

func NewEpisode() *Episode {
	return &Episode{}
}
