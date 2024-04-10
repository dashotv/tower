// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/dashotv/flame/qbt"
	"github.com/dashotv/grimoire"
)

func init() {
	initializers = append(initializers, setupDb)
	healthchecks["db"] = checkDb
}

func setupDb(app *Application) error {
	db, err := NewConnector(app)
	if err != nil {
		return err
	}

	app.DB = db
	return nil
}

func checkDb(app *Application) (err error) {
	// TODO: Check DB connection
	return nil
}

type Connector struct {
	Log         *zap.SugaredLogger
	Collection  *grimoire.Store[*Collection]
	Combination *grimoire.Store[*Combination]
	Download    *grimoire.Store[*Download]
	Episode     *grimoire.Store[*Episode]
	Feed        *grimoire.Store[*Feed]
	File        *grimoire.Store[*File]
	Medium      *grimoire.Store[*Medium]
	Message     *grimoire.Store[*Message]
	Minion      *grimoire.Store[*Minion]
	Movie       *grimoire.Store[*Movie]
	Pin         *grimoire.Store[*Pin]
	Release     *grimoire.Store[*Release]
	Request     *grimoire.Store[*Request]
	Series      *grimoire.Store[*Series]
	User        *grimoire.Store[*User]
	Watch       *grimoire.Store[*Watch]
}

func NewConnector(app *Application) (*Connector, error) {
	var s *Connection
	var err error

	s, err = app.Config.ConnectionFor("collection")
	if err != nil {
		return nil, err
	}
	collection, err := grimoire.New[*Collection](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Collection](collection, &Collection{})

	s, err = app.Config.ConnectionFor("combination")
	if err != nil {
		return nil, err
	}
	combination, err := grimoire.New[*Combination](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Combination](combination, &Combination{})

	s, err = app.Config.ConnectionFor("download")
	if err != nil {
		return nil, err
	}
	download, err := grimoire.New[*Download](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Download](download, &Download{})

	s, err = app.Config.ConnectionFor("episode")
	if err != nil {
		return nil, err
	}
	episode, err := grimoire.New[*Episode](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Episode](episode, &Episode{})

	s, err = app.Config.ConnectionFor("feed")
	if err != nil {
		return nil, err
	}
	feed, err := grimoire.New[*Feed](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Feed](feed, &Feed{})

	s, err = app.Config.ConnectionFor("file")
	if err != nil {
		return nil, err
	}
	file, err := grimoire.New[*File](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*File](file, &File{})

	s, err = app.Config.ConnectionFor("medium")
	if err != nil {
		return nil, err
	}
	medium, err := grimoire.New[*Medium](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Medium](medium, &Medium{})

	s, err = app.Config.ConnectionFor("message")
	if err != nil {
		return nil, err
	}
	message, err := grimoire.New[*Message](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Message](message, &Message{})

	s, err = app.Config.ConnectionFor("minion")
	if err != nil {
		return nil, err
	}
	minion, err := grimoire.New[*Minion](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Minion](minion, &Minion{})

	s, err = app.Config.ConnectionFor("movie")
	if err != nil {
		return nil, err
	}
	movie, err := grimoire.New[*Movie](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Movie](movie, &Movie{})

	s, err = app.Config.ConnectionFor("pin")
	if err != nil {
		return nil, err
	}
	pin, err := grimoire.New[*Pin](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Pin](pin, &Pin{})

	s, err = app.Config.ConnectionFor("release")
	if err != nil {
		return nil, err
	}
	release, err := grimoire.New[*Release](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Release](release, &Release{})

	s, err = app.Config.ConnectionFor("request")
	if err != nil {
		return nil, err
	}
	request, err := grimoire.New[*Request](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Request](request, &Request{})

	s, err = app.Config.ConnectionFor("series")
	if err != nil {
		return nil, err
	}
	series, err := grimoire.New[*Series](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Series](series, &Series{})

	s, err = app.Config.ConnectionFor("user")
	if err != nil {
		return nil, err
	}
	user, err := grimoire.New[*User](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*User](user, &User{})

	s, err = app.Config.ConnectionFor("watch")
	if err != nil {
		return nil, err
	}
	watch, err := grimoire.New[*Watch](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	grimoire.Indexes[*Watch](watch, &Watch{})

	c := &Connector{
		Log:         app.Log.Named("db"),
		Collection:  collection,
		Combination: combination,
		Download:    download,
		Episode:     episode,
		Feed:        feed,
		File:        file,
		Medium:      medium,
		Message:     message,
		Minion:      minion,
		Movie:       movie,
		Pin:         pin,
		Release:     release,
		Request:     request,
		Series:      series,
		User:        user,
		Watch:       watch,
	}

	return c, nil
}

type Collection struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name      string             `bson:"name" json:"name"`
	Library   string             `bson:"library" json:"library"`
	RatingKey string             `bson:"rating_key" json:"rating_key"`
	SyncedAt  time.Time          `bson:"synced_at" json:"synced_at"`
	Media     []*CollectionMedia `bson:"media" json:"media"`
}

type CollectionMedia struct { // struct
	RatingKey string `bson:"rating_key" json:"rating_key"`
	Title     string `bson:"title" json:"title"`
}

type Combination struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name        string   `bson:"name" json:"name"`
	Collections []string `bson:"collections" json:"collections"`
}

type CombinationChild struct { // struct
	RatingKey    string `bson:"rating_key" json:"rating_key"`
	Key          string `bson:"key" json:"key"`
	GUID         string `bson:"guid" json:"guid"`
	Type         string `bson:"type" json:"type"`
	Title        string `bson:"title" json:"title"`
	LibraryID    int64  `bson:"library_id" json:"library_id"`
	LibraryTitle string `bson:"library_title" json:"library_title"`
	LibraryKey   string `bson:"library_key" json:"library_key"`
	Summary      string `bson:"summary" json:"summary"`
	Thumb        string `bson:"thumb" json:"thumb"`
	Total        int    `bson:"total" json:"total"`
	Viewed       int    `bson:"viewed" json:"viewed"`
	Link         string `bson:"link" json:"link"`
	Next         string `bson:"next" json:"next"`
	LastViewedAt int64  `bson:"last_viewed_at" json:"last_viewed_at"`
	AddedAt      int64  `bson:"added_at" json:"added_at"`
	UpdatedAt    int64  `bson:"updated_at" json:"updated_at"`
}

type Download struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	MediumID     primitive.ObjectID `bson:"medium_id" json:"medium_id"`
	Auto         bool               `bson:"auto" json:"auto"`
	Multi        bool               `bson:"multi" json:"multi"`
	Force        bool               `bson:"force" json:"force"`
	URL          string             `bson:"url" json:"url"`
	ReleaseID    string             `bson:"tdo_id" json:"release_id"`
	Thash        string             `bson:"thash" json:"thash"`
	Selected     string             `bson:"selected" json:"selected"`
	Status       string             `bson:"status" json:"status"`
	Files        []*DownloadFile    `bson:"download_files" json:"files"`
	Medium       *Medium            `bson:"-" json:"medium"`
	Title        string             `bson:"-" json:"title"`
	Display      string             `bson:"-" json:"display"`
	Source       string             `bson:"-" json:"source"`
	SourceID     string             `bson:"-" json:"source_id"`
	Kind         primitive.Symbol   `bson:"-" json:"kind"`
	Search       string             `bson:"-" json:"search"`
	Directory    string             `bson:"-" json:"directory"`
	SearchParams *SearchParams      `bson:"-" json:"search_params"`
	Active       bool               `bson:"-" json:"active"`
	Favorite     bool               `bson:"-" json:"favorite"`
	Unwatched    int                `bson:"-" json:"unwatched"`
	Cover        string             `bson:"-" json:"cover"`
	Background   string             `bson:"-" json:"background"`
}

type DownloadFile struct { // struct
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	MediumID    primitive.ObjectID `bson:"medium_id" json:"medium_id"`
	Medium      *Medium            `bson:"medium" json:"medium"`
	Num         int                `bson:"num" json:"num"`
	TorrentFile *TorrentFile       `bson:"-" json:"-"`
}

type Downloading struct { // struct
	ID           string            `bson:"id" json:"id"`
	MediumID     string            `bson:"medium_id" json:"medium_id"`
	Multi        bool              `bson:"multi" json:"multi"`
	Infohash     string            `bson:"infohash" json:"infohash"`
	Torrent      *qbt.TorrentJSON  `bson:"torrent" json:"torrent"`
	Queue        float64           `bson:"queue" json:"queue"`
	Progress     float64           `bson:"progress" json:"progress"`
	Eta          string            `bson:"eta" json:"eta"`
	TorrentState string            `bson:"torrent_state" json:"torrent_state"`
	Files        *DownloadingFiles `bson:"files" json:"files"`
	Title        string            `bson:"title" json:"title"`
	Display      string            `bson:"display" json:"display"`
	Cover        string            `bson:"cover" json:"cover"`
	Background   string            `bson:"background" json:"background"`
}

type DownloadingFiles struct { // struct
	Completed int `bson:"completed" json:"completed"`
	Selected  int `bson:"selected" json:"selected"`
}

type Episode struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Type            string             `bson:"_type" json:"type"`
	SeriesID        primitive.ObjectID `bson:"series_id" json:"series_id"`
	SourceID        string             `bson:"source_id" json:"source_id"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Directory       string             `bson:"directory" json:"directory"`
	Search          string             `bson:"search" json:"search"`
	SeasonNumber    int                `bson:"season_number" json:"season_number"`
	EpisodeNumber   int                `bson:"episode_number" json:"episode_number"`
	AbsoluteNumber  int                `bson:"absolute_number" json:"absolute_number"`
	Downloaded      bool               `bson:"downloaded" json:"downloaded"`
	Completed       bool               `bson:"completed" json:"completed"`
	Skipped         bool               `bson:"skipped" json:"skipped"`
	Missing         *time.Time         `bson:"missing" json:"missing"`
	ReleaseDate     time.Time          `bson:"release_date" json:"release_date"`
	Paths           []*Path            `bson:"paths,omitempty" json:"paths"`
	Cover           string             `bson:"-" json:"cover"`
	Background      string             `bson:"-" json:"background"`
	Watched         bool               `bson:"-" json:"watched"`
	WatchedAny      bool               `bson:"-" json:"watched_any"`
	SeriesTitle     string             `bson:"-" json:"series_title"`
	SeriesDisplay   string             `bson:"-" json:"series_display"`
	SeriesSource    string             `bson:"-" json:"series_source"`
	SeriesKind      primitive.Symbol   `bson:"-" json:"series_kind"`
	SeriesActive    bool               `bson:"-" json:"series_active"`
	SeriesFavorite  bool               `bson:"-" json:"series_favorite"`
	SeriesUnwatched int                `bson:"-" json:"series_unwatched"`
}

type Feed struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name      string    `bson:"name" json:"name"`
	URL       string    `bson:"url" json:"url"`
	Source    string    `bson:"source" json:"source"`
	Type      string    `bson:"type" json:"type"`
	Active    bool      `bson:"active" json:"active"`
	Processed time.Time `bson:"processed" json:"processed"`
}

type File struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Type       string             `bson:"type" json:"type"`
	Path       string             `bson:"path" json:"path"`
	Size       int64              `bson:"size" json:"size"`
	ModifiedAt int64              `bson:"modified_at" json:"modified_at"`
	MediumID   primitive.ObjectID `bson:"medium_id" json:"medium_id"`
}

type Medium struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Type           string             `bson:"_type" json:"type"`
	Kind           primitive.Symbol   `bson:"kind" json:"kind"`
	Source         string             `bson:"source" json:"source"`
	SourceID       string             `bson:"source_id" json:"source_id"`
	ImdbID         string             `bson:"imdb_id" json:"imdb_id"`
	Title          string             `bson:"title" json:"title"`
	Description    string             `bson:"description" json:"description"`
	Display        string             `bson:"display" json:"display"`
	Directory      string             `bson:"directory" json:"directory"`
	Search         string             `bson:"search" json:"search"`
	SearchParams   *SearchParams      `bson:"search_params" json:"search_params"`
	Active         bool               `bson:"active" json:"active"`
	Downloaded     bool               `bson:"downloaded" json:"downloaded"`
	Completed      bool               `bson:"completed" json:"completed"`
	Skipped        bool               `bson:"skipped" json:"skipped"`
	Watched        bool               `bson:"watched" json:"watched"`
	Broken         bool               `bson:"broken" json:"broken"`
	Favorite       bool               `bson:"favorite" json:"favorite"`
	Unwatched      int                `bson:"unwatched" json:"unwatched"`
	ReleaseDate    time.Time          `bson:"release_date" json:"release_date"`
	Paths          []*Path            `bson:"paths,omitempty" json:"paths"`
	Cover          string             `bson:"-" json:"cover"`
	Background     string             `bson:"-" json:"background"`
	SeriesID       primitive.ObjectID `bson:"series_id" json:"series_id"`
	SeasonNumber   int                `bson:"season_number" json:"season_number"`
	EpisodeNumber  int                `bson:"episode_number" json:"episode_number"`
	AbsoluteNumber int                `bson:"absolute_number" json:"absolute_number"`
}

type Message struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Level    string `bson:"level" json:"level"`
	Facility string `bson:"facility" json:"facility"`
	Message  string `bson:"message" json:"message"`
}

type Minion struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Kind     string           `bson:"kind" json:"kind"`
	Args     string           `bson:"args" json:"args"`
	Status   string           `bson:"status" json:"status"`
	Queue    string           `bson:"queue" json:"queue"`
	Attempts []*MinionAttempt `bson:"attempts" json:"attempts"`
}

type MinionAttempt struct { // struct
	StartedAt  time.Time `bson:"started_at" json:"started_at"`
	Duration   float64   `bson:"duration" json:"duration"`
	Status     string    `bson:"status" json:"status"`
	Error      string    `bson:"error" json:"error"`
	Stacktrace []string  `bson:"stacktrace" json:"stacktrace"`
}

type Movie struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Type         string           `bson:"_type" json:"type"`
	Kind         primitive.Symbol `bson:"kind" json:"kind"`
	Source       string           `bson:"source" json:"source"`
	SourceID     string           `bson:"source_id" json:"source_id"`
	ImdbID       string           `bson:"imdb_id" json:"imdb_id"`
	Title        string           `bson:"title" json:"title"`
	Description  string           `bson:"description" json:"description"`
	Slug         string           `bson:"slug" json:"slug"`
	Text         []string         `bson:"text" json:"text"`
	Display      string           `bson:"display" json:"display"`
	Directory    string           `bson:"directory" json:"directory"`
	Search       string           `bson:"search" json:"search"`
	SearchParams *SearchParams    `bson:"search_params" json:"search_params"`
	Active       bool             `bson:"active" json:"active"`
	Downloaded   bool             `bson:"downloaded" json:"downloaded"`
	Completed    bool             `bson:"completed" json:"completed"`
	Skipped      bool             `bson:"skipped" json:"skipped"`
	Watched      bool             `bson:"watched" json:"watched"`
	Broken       bool             `bson:"broken" json:"broken"`
	Favorite     bool             `bson:"favorite" json:"favorite"`
	ReleaseDate  time.Time        `bson:"release_date" json:"release_date"`
	Paths        []*Path          `bson:"paths,omitempty" json:"paths"`
	Cover        string           `bson:"-" json:"cover"`
	Background   string           `bson:"-" json:"background"`
}

type NzbgetPayload struct { // struct
	ID       string `bson:"id" json:"id"`
	Name     string `bson:"name" json:"name"`
	Category string `bson:"category" json:"category"`
	Dir      string `bson:"dir" json:"dir"`
	FinalDir string `bson:"final_dir" json:"final_dir"`
	File     string `bson:"file" json:"file"`
	Status   string `bson:"status" json:"status"`
}

type Path struct { // struct
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type       primitive.Symbol   `bson:"type" json:"type"`
	Remote     string             `bson:"remote" json:"remote"`
	Local      string             `bson:"local" json:"local"`
	Extension  string             `bson:"extension" json:"extension"`
	Size       int                `bson:"size" json:"size"`
	Resolution int                `bson:"resolution" json:"resolution"`
	Bitrate    int                `bson:"bitrate" json:"bitrate"`
	Checksum   string             `bson:"checksum" json:"checksum"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type Pin struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Pin        int    `bson:"pin" json:"id"`
	Code       string `bson:"code" json:"code"`
	Token      string `bson:"token" json:"authToken"`
	Product    string `bson:"product" json:"product"`
	Identifier string `bson:"identifier" json:"clientIdentifier"`
}

type Popular struct { // struct
	Name  string `bson:"_id" json:"name"`
	Year  int    `bson:"year" json:"year"`
	Type  string `bson:"type" json:"type"`
	Count int    `bson:"count" json:"count"`
}

type PopularResponse struct { // struct
	Tv     []*Popular `bson:"tv" json:"tv"`
	Anime  []*Popular `bson:"anime" json:"anime"`
	Movies []*Popular `bson:"movies" json:"movies"`
}

type Release struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Type        string    `bson:"type" json:"type"`
	Source      string    `bson:"source" json:"source"`
	Raw         string    `bson:"raw" json:"raw"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	Size        string    `bson:"size" json:"size"`
	View        string    `bson:"view" json:"view"`
	Download    string    `bson:"download" json:"download"`
	Infohash    string    `bson:"infohash" json:"infohash"`
	Name        string    `bson:"name" json:"name"`
	Season      int       `bson:"season" json:"season"`
	Episode     int       `bson:"episode" json:"episode"`
	Volume      int       `bson:"volume" json:"volume"`
	Year        int       `bson:"year" json:"year"`
	Checksum    string    `bson:"checksum" json:"checksum"`
	Group       string    `bson:"group" json:"group"`
	Author      string    `bson:"author" json:"author"`
	Verified    bool      `bson:"verified" json:"verified"`
	Widescreen  bool      `bson:"widescreen" json:"widescreen"`
	Uncensored  bool      `bson:"uncensored" json:"uncensored"`
	Bluray      bool      `bson:"bluray" json:"bluray"`
	Nzb         bool      `bson:"nzb" json:"nzb"`
	Resolution  string    `bson:"resolution" json:"resolution"`
	Encoding    string    `bson:"encoding" json:"encoding"`
	Quality     string    `bson:"quality" json:"quality"`
	PublishedAt time.Time `bson:"published_at" json:"published_at"`
}

type Request struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Title    string `bson:"title" json:"title"`
	User     string `bson:"user" json:"user"`
	Type     string `bson:"type" json:"type"`
	Source   string `bson:"source" json:"source"`
	SourceID string `bson:"source_id" json:"source_id"`
	Status   string `bson:"status" json:"status"`
}

type SearchParams struct { // struct
	Type       string `bson:"type" json:"type"`
	Verified   bool   `bson:"verified" json:"verified"`
	Group      string `bson:"group" json:"group"`
	Author     string `bson:"author" json:"author"`
	Resolution int    `bson:"resolution" json:"resolution"`
	Source     string `bson:"source" json:"source"`
	Uncensored bool   `bson:"uncensored" json:"uncensored"`
	Bluray     bool   `bson:"bluray" json:"bluray"`
}

type Series struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Type          string           `bson:"_type" json:"type"`
	Kind          primitive.Symbol `bson:"kind" json:"kind"`
	Source        string           `bson:"source" json:"source"`
	SourceID      string           `bson:"source_id" json:"source_id"`
	ImdbID        string           `bson:"imdb_id" json:"imdb_id"`
	Title         string           `bson:"title" json:"title"`
	Description   string           `bson:"description" json:"description"`
	Slug          string           `bson:"slug" json:"slug"`
	Text          []string         `bson:"text" json:"text"`
	Display       string           `bson:"display" json:"display"`
	Directory     string           `bson:"directory" json:"directory"`
	Search        string           `bson:"search" json:"search"`
	SearchParams  *SearchParams    `bson:"search_params" json:"search_params"`
	Status        string           `bson:"status" json:"status"`
	Active        bool             `bson:"active" json:"active"`
	Downloaded    bool             `bson:"downloaded" json:"downloaded"`
	Completed     bool             `bson:"completed" json:"completed"`
	Skipped       bool             `bson:"skipped" json:"skipped"`
	Watched       bool             `bson:"watched" json:"watched"`
	Broken        bool             `bson:"broken" json:"broken"`
	Favorite      bool             `bson:"favorite" json:"favorite"`
	Unwatched     int              `bson:"-" json:"unwatched"`
	UnwatchedAll  int              `bson:"-" json:"unwatched_all"`
	ReleaseDate   time.Time        `bson:"release_date" json:"release_date"`
	Paths         []*Path          `bson:"paths,omitempty" json:"paths"`
	Cover         string           `bson:"-" json:"cover"`
	Background    string           `bson:"-" json:"background"`
	CurrentSeason int              `bson:"-" json:"currentSeason"`
	Seasons       []int            `bson:"-" json:"seasons"`
	Episodes      []*Episode       `bson:"-" json:"episodes"`
	Watches       []*Watch         `bson:"-" json:"watches"`
}

type TorrentFile struct { // struct
	ID       int     `bson:"id" json:"id"`
	IsSend   bool    `bson:"is_send" json:"is_send"`
	Name     string  `bson:"name" json:"name"`
	Priority int     `bson:"priority" json:"priority"`
	Progress float64 `bson:"progress" json:"progress"`
	Size     int64   `bson:"size" json:"size"`
}

type Upcoming struct { // struct
	ID                      primitive.ObjectID `bson:"id" json:"id"`
	Type                    string             `bson:"type" json:"type"`
	SourceID                string             `bson:"source_id" json:"source_id"`
	Title                   string             `bson:"title" json:"title"`
	Display                 string             `bson:"display" json:"display"`
	Description             string             `bson:"description" json:"description"`
	Directory               string             `bson:"directory" json:"directory"`
	Search                  string             `bson:"search" json:"search"`
	SeasonNumber            int                `bson:"season_number" json:"season_number"`
	EpisodeNumber           int                `bson:"episode_number" json:"episode_number"`
	AbsoluteNumber          int                `bson:"absolute_number" json:"absolute_number"`
	Downloaded              bool               `bson:"downloaded" json:"downloaded"`
	Completed               bool               `bson:"completed" json:"completed"`
	Skipped                 bool               `bson:"skipped" json:"skipped"`
	ReleaseDate             time.Time          `bson:"release_date" json:"release_date"`
	SeriesID                primitive.ObjectID `bson:"series_id" json:"series_id"`
	SeriesSource            string             `bson:"series_source" json:"series_source"`
	SeriesTitle             string             `bson:"series_title" json:"series_title"`
	SeriesKind              primitive.Symbol   `bson:"series_kind" json:"series_kind"`
	SeriesActive            bool               `bson:"series_active" json:"series_active"`
	SeriesFavorite          bool               `bson:"series_favorite" json:"series_favorite"`
	SeriesUnwatched         int                `bson:"series_unwatched" json:"series_unwatched"`
	SeriesUnwatchedAll      int                `bson:"series_unwatched_all" json:"series_unwatched_all"`
	SeriesCover             string             `bson:"series_cover" json:"series_cover"`
	SeriesCoverUpdated      time.Time          `bson:"series_cover_updated" json:"series_cover_updated"`
	SeriesBackground        string             `bson:"series_background" json:"series_background"`
	SeriesBackgroundUpdated time.Time          `bson:"series_background_updated" json:"series_background_updated"`
}

type User struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Token string `bson:"token" json:"token"`
	Thumb string `bson:"thumb" json:"thumb"`
	Home  bool   `bson:"home" json:"home"`
	Admin bool   `bson:"admin" json:"admin"`
}

type Watch struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Username  string             `bson:"username" json:"username"`
	Player    string             `bson:"player" json:"player"`
	WatchedAt time.Time          `bson:"watched_at" json:"watched_at"`
	MediumID  primitive.ObjectID `bson:"medium_id" json:"medium_id"`
	Medium    *Medium            `bson:"-" json:"medium"`
}
