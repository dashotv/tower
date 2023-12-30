// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

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
	Log      *zap.SugaredLogger
	Download *grimoire.Store[*Download]
	Episode  *grimoire.Store[*Episode]
	Feed     *grimoire.Store[*Feed]
	Medium   *grimoire.Store[*Medium]
	Message  *grimoire.Store[*Message]
	Minion   *grimoire.Store[*Minion]
	Movie    *grimoire.Store[*Movie]
	Pin      *grimoire.Store[*Pin]
	Release  *grimoire.Store[*Release]
	Request  *grimoire.Store[*Request]
	Series   *grimoire.Store[*Series]
	User     *grimoire.Store[*User]
	Watch    *grimoire.Store[*Watch]
}

func NewConnector(app *Application) (*Connector, error) {
	var s *Connection
	var err error

	s, err = app.Config.ConnectionFor("download")
	if err != nil {
		return nil, err
	}
	download, err := grimoire.New[*Download](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("episode")
	if err != nil {
		return nil, err
	}
	episode, err := grimoire.New[*Episode](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("feed")
	if err != nil {
		return nil, err
	}
	feed, err := grimoire.New[*Feed](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("medium")
	if err != nil {
		return nil, err
	}
	medium, err := grimoire.New[*Medium](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("message")
	if err != nil {
		return nil, err
	}
	message, err := grimoire.New[*Message](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("minion")
	if err != nil {
		return nil, err
	}
	minion, err := grimoire.New[*Minion](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("movie")
	if err != nil {
		return nil, err
	}
	movie, err := grimoire.New[*Movie](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("pin")
	if err != nil {
		return nil, err
	}
	pin, err := grimoire.New[*Pin](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("release")
	if err != nil {
		return nil, err
	}
	release, err := grimoire.New[*Release](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("request")
	if err != nil {
		return nil, err
	}
	request, err := grimoire.New[*Request](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("series")
	if err != nil {
		return nil, err
	}
	series, err := grimoire.New[*Series](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("user")
	if err != nil {
		return nil, err
	}
	user, err := grimoire.New[*User](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = app.Config.ConnectionFor("watch")
	if err != nil {
		return nil, err
	}
	watch, err := grimoire.New[*Watch](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	c := &Connector{
		Log:      app.Log.Named("db"),
		Download: download,
		Episode:  episode,
		Feed:     feed,
		Medium:   medium,
		Message:  message,
		Minion:   minion,
		Movie:    movie,
		Pin:      pin,
		Release:  release,
		Request:  request,
		Series:   series,
		User:     user,
		Watch:    watch,
	}

	return c, nil
}

type Download struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	MediumId  primitive.ObjectID `bson:"medium_id" json:"medium_id"`
	Auto      bool               `bson:"auto" json:"auto"`
	Multi     bool               `bson:"multi" json:"multi"`
	Force     bool               `bson:"force" json:"force"`
	Url       string             `bson:"url" json:"url"`
	ReleaseId string             `bson:"tdo_id" json:"release_id"`
	Thash     string             `bson:"thash" json:"thash"`
	Selected  string             `bson:"selected" json:"selected"`
	Status    string             `bson:"status" json:"status"`
	Files     []*DownloadFile    `bson:"download_files" json:"download_files"`
	Medium    *Medium            `bson:"-" json:"medium"`
}

type DownloadFile struct { // struct
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	MediumId primitive.ObjectID `bson:"medium_id" json:"medium_id"`
	Medium   *Medium            `bson:"medium" json:"medium"`
	Num      int                `bson:"num" json:"num"`
}

type Episode struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Type           string             `bson:"_type" json:"type"`
	SeriesId       primitive.ObjectID `bson:"series_id" json:"series_id"`
	SourceId       string             `bson:"source_id" json:"source_id"`
	Title          string             `bson:"title" json:"title"`
	Description    string             `bson:"description" json:"description"`
	Directory      string             `bson:"directory" json:"directory"`
	Search         string             `bson:"search" json:"search"`
	SeasonNumber   int                `bson:"season_number" json:"season_number"`
	EpisodeNumber  int                `bson:"episode_number" json:"episode_number"`
	AbsoluteNumber int                `bson:"absolute_number" json:"absolute_number"`
	Downloaded     bool               `bson:"downloaded" json:"downloaded"`
	Completed      bool               `bson:"completed" json:"completed"`
	Skipped        bool               `bson:"skipped" json:"skipped"`
	ReleaseDate    time.Time          `bson:"release_date" json:"release_date"`
	Paths          []*Path            `bson:"paths,omitempty" json:"paths"`
	Cover          string             `bson:"-" json:"cover"`
	Background     string             `bson:"-" json:"background"`
	Watched        bool               `bson:"-" json:"watched"`
	Active         bool               `bson:"-" json:"active"`
	Favorite       bool               `bson:"-" json:"favorite"`
	Unwatched      int                `bson:"-" json:"unwatched"`
	Display        string             `bson:"-" json:"display"`
}

type Feed struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name      string    `bson:"name" json:"name"`
	Url       string    `bson:"url" json:"url"`
	Source    string    `bson:"source" json:"source"`
	Type      string    `bson:"type" json:"type"`
	Active    bool      `bson:"active" json:"active"`
	Processed time.Time `bson:"processed" json:"processed"`
}

type Medium struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Type           string             `bson:"_type" json:"type"`
	Kind           primitive.Symbol   `bson:"kind" json:"kind"`
	Source         string             `bson:"source" json:"source"`
	SourceId       string             `bson:"source_id" json:"source_id"`
	ImdbId         string             `bson:"imdb_id" json:"imdb_id"`
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
	SeriesId       primitive.ObjectID `bson:"series_id" json:"series_id"`
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
	SourceId     string           `bson:"source_id" json:"source_id"`
	ImdbId       string           `bson:"imdb_id" json:"imdb_id"`
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

type Path struct { // struct
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
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
	SourceId string `bson:"source_id" json:"source_id"`
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
	SourceId      string           `bson:"source_id" json:"source_id"`
	ImdbId        string           `bson:"imdb_id" json:"imdb_id"`
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
	ReleaseDate   time.Time        `bson:"release_date" json:"release_date"`
	Paths         []*Path          `bson:"paths,omitempty" json:"paths"`
	Cover         string           `bson:"-" json:"cover"`
	Background    string           `bson:"-" json:"background"`
	CurrentSeason int              `bson:"-" json:"currentSeason"`
	Seasons       []int            `bson:"-" json:"seasons"`
	Episodes      []*Episode       `bson:"-" json:"episodes"`
	Watches       []*Watch         `bson:"-" json:"watches"`
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
	MediumId  primitive.ObjectID `bson:"medium_id" json:"medium_id"`
	Medium    *Medium            `bson:"-" json:"medium"`
}