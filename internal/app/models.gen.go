// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/dashotv/flame/qbt"
	"github.com/dashotv/grimoire"
	"github.com/kamva/mgm/v3"
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
	Log             *zap.SugaredLogger
	Collection      *grimoire.Store[*Collection]
	Combination     *grimoire.Store[*Combination]
	Download        *grimoire.Store[*Download]
	Episode         *grimoire.Store[*Episode]
	Feed            *grimoire.Store[*Feed]
	File            *grimoire.Store[*File]
	Library         *grimoire.Store[*Library]
	LibraryTemplate *grimoire.Store[*LibraryTemplate]
	LibraryType     *grimoire.Store[*LibraryType]
	Medium          *grimoire.Store[*Medium]
	Message         *grimoire.Store[*Message]
	Movie           *grimoire.Store[*Movie]
	Pin             *grimoire.Store[*Pin]
	Release         *grimoire.Store[*Release]
	Request         *grimoire.Store[*Request]
	Series          *grimoire.Store[*Series]
	User            *grimoire.Store[*User]
	Watch           *grimoire.Store[*Watch]
}

func connection[T mgm.Model](name string) (*grimoire.Store[T], error) {
	s, err := app.Config.ConnectionFor(name)
	if err != nil {
		return nil, err
	}
	c, err := grimoire.New[T](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewConnector(app *Application) (*Connector, error) {
	collection, err := connection[*Collection]("collection")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*Collection](collection, &Collection{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*Collection](collection, &Collection{})

	combination, err := connection[*Combination]("combination")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*Combination](combination, &Combination{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*Combination](combination, &Combination{})

	download, err := connection[*Download]("download")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*Download](download, &Download{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*Download](download, &Download{})

	episode, err := connection[*Episode]("episode")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexesFromTags[*Episode](episode, &Episode{})

	feed, err := connection[*Feed]("feed")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexesFromTags[*Feed](feed, &Feed{})

	file, err := connection[*File]("file")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*File](file, &File{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*File](file, &File{})

	library, err := connection[*Library]("library")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*Library](library, &Library{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*Library](library, &Library{})

	library_template, err := connection[*LibraryTemplate]("library_template")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*LibraryTemplate](library_template, &LibraryTemplate{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*LibraryTemplate](library_template, &LibraryTemplate{})

	library_type, err := connection[*LibraryType]("library_type")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*LibraryType](library_type, &LibraryType{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*LibraryType](library_type, &LibraryType{})

	medium, err := connection[*Medium]("medium")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexesFromTags[*Medium](medium, &Medium{})

	message, err := connection[*Message]("message")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*Message](message, &Message{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*Message](message, &Message{})

	movie, err := connection[*Movie]("movie")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexesFromTags[*Movie](movie, &Movie{})

	pin, err := connection[*Pin]("pin")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*Pin](pin, &Pin{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*Pin](pin, &Pin{})

	release, err := connection[*Release]("release")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexesFromTags[*Release](release, &Release{})

	request, err := connection[*Request]("request")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*Request](request, &Request{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*Request](request, &Request{})

	series, err := connection[*Series]("series")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexesFromTags[*Series](series, &Series{})

	user, err := connection[*User]("user")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*User](user, &User{}, "created_at;updated_at")

	grimoire.CreateIndexesFromTags[*User](user, &User{})

	watch, err := connection[*Watch]("watch")
	if err != nil {
		return nil, err
	}

	grimoire.CreateIndexes[*Watch](watch, &Watch{}, "created_at;updated_at;username;medium_id;watched_at")

	grimoire.CreateIndexesFromTags[*Watch](watch, &Watch{})

	c := &Connector{
		Log:             app.Log.Named("db"),
		Collection:      collection,
		Combination:     combination,
		Download:        download,
		Episode:         episode,
		Feed:            feed,
		File:            file,
		Library:         library,
		LibraryTemplate: library_template,
		LibraryType:     library_type,
		Medium:          medium,
		Message:         message,
		Movie:           movie,
		Pin:             pin,
		Release:         release,
		Request:         request,
		Series:          series,
		User:            user,
		Watch:           watch,
	}

	return c, nil
}

type Collection struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name      string             `bson:"name" json:"name" grimoire:"index"`
	Library   string             `bson:"library" json:"library" grimoire:"index"`
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
	Name        string   `bson:"name" json:"name" grimoire:"index"`
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

type Directory struct { // struct
	Name  string `bson:"name" json:"name"`
	Path  string `bson:"path" json:"path"`
	Count int64  `bson:"count" json:"count"`
}

type Download struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	MediumID       primitive.ObjectID `bson:"medium_id" json:"medium_id"`
	Auto           bool               `bson:"auto" json:"auto"`
	Multi          bool               `bson:"multi" json:"multi"`
	Force          bool               `bson:"force" json:"force"`
	URL            string             `bson:"url" json:"url"`
	ReleaseID      string             `bson:"tdo_id" json:"release_id"`
	Thash          string             `bson:"thash" json:"thash"`
	Selected       string             `bson:"selected" json:"selected"`
	Status         string             `bson:"status" json:"status"`
	Files          []*DownloadFile    `bson:"download_files" json:"files"`
	Regex          string             `bson:"regex" json:"regex"`
	Tag            string             `bson:"tag" json:"tag"`
	Medium         *Medium            `bson:"-" json:"medium"`
	Title          string             `bson:"-" json:"title"`
	Display        string             `bson:"-" json:"display"`
	Source         string             `bson:"-" json:"source"`
	SourceID       string             `bson:"-" json:"source_id"`
	Kind           primitive.Symbol   `bson:"-" json:"kind"`
	Directory      string             `bson:"-" json:"directory"`
	Active         bool               `bson:"-" json:"active"`
	Favorite       bool               `bson:"-" json:"favorite"`
	Unwatched      int                `bson:"-" json:"unwatched"`
	Cover          string             `bson:"-" json:"cover"`
	Background     string             `bson:"-" json:"background"`
	Search         *DownloadSearch    `bson:"-" json:"search"`
	Torrent        *qbt.TorrentJSON   `bson:"-" json:"torrent"`
	TorrentState   string             `bson:"-" json:"torrent_state"`
	Eta            string             `bson:"-" json:"eta"`
	Progress       float64            `bson:"-" json:"progress"`
	Queue          float64            `bson:"-" json:"queue"`
	FilesCompleted int                `bson:"-" json:"files_completed"`
	FilesSelected  int                `bson:"-" json:"files_selected"`
	FilesWanted    int                `bson:"-" json:"files_wanted"`
}

type DownloadFile struct { // struct
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	MediumID    primitive.ObjectID `bson:"medium_id" json:"medium_id"`
	Num         int                `bson:"num" json:"num"`
	Medium      *Medium            `bson:"-" json:"medium"`
	TorrentFile *qbt.TorrentFile   `bson:"-" json:"torrent_file"`
}

type DownloadSearch struct { // struct
	Type       string `bson:"type" json:"type"`
	Source     string `bson:"source" json:"source"`
	SourceID   string `bson:"source_id" json:"source_id"`
	Title      string `bson:"title" json:"title"`
	Year       int    `bson:"year" json:"year"`
	Season     int    `bson:"season" json:"season"`
	Episode    int    `bson:"episode" json:"episode"`
	Resolution int    `bson:"resolution" json:"resolution"`
	Group      string `bson:"group" json:"group"`
	Website    string `bson:"website" json:"website"`
	Exact      bool   `bson:"exact" json:"exact"`
	Verified   bool   `bson:"verified" json:"verified"`
	Uncensored bool   `bson:"uncensored" json:"uncensored"`
	Bluray     bool   `bson:"bluray" json:"bluray"`
	ThreeD     bool   `bson:"three_d" json:"three_d"`
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
	SeasonNumber    int                `bson:"season_number" json:"season_number"`
	EpisodeNumber   int                `bson:"episode_number" json:"episode_number"`
	AbsoluteNumber  int                `bson:"absolute_number" json:"absolute_number"`
	Downloaded      bool               `bson:"downloaded" json:"downloaded"`
	Completed       bool               `bson:"completed" json:"completed"`
	Skipped         bool               `bson:"skipped" json:"skipped"`
	Missing         *time.Time         `bson:"missing" json:"missing"`
	ReleaseDate     time.Time          `bson:"release_date" json:"release_date"`
	Overrides       *Overrides         `bson:"overrides" json:"overrides"`
	Paths           []*Path            `bson:"paths,omitempty" json:"paths"`
	Cover           string             `bson:"-" json:"cover"`
	Background      string             `bson:"-" json:"background"`
	HasOverrides    bool               `bson:"-" json:"has_overrides"`
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
	LibraryID  primitive.ObjectID `bson:"library_id" json:"library_id" grimoire:"index"`
	MediumID   primitive.ObjectID `bson:"medium_id" json:"medium_id" grimoire:"index"`
	Type       string             `bson:"type" json:"type" grimoire:"index"`
	Name       string             `bson:"name" json:"name"`
	Extension  string             `bson:"extension" json:"extension"`
	Path       string             `bson:"path" json:"path" grimoire:"index"`
	Size       int64              `bson:"size" json:"size"`
	Resolution int                `bson:"resolution" json:"resolution"`
	Checksum   string             `bson:"checksum" json:"checksum"`
	ModifiedAt int64              `bson:"modified_at" json:"modified_at"`
	Exists     bool               `bson:"exists" json:"exists"`
	Old        bool               `bson:"old" json:"old"`
	Rename     bool               `bson:"rename" json:"rename"`
}

type Library struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name              string             `bson:"name" json:"name" grimoire:"index"`
	Path              string             `bson:"path" json:"path" grimoire:"index"`
	Count             int64              `bson:"count" json:"count"`
	LibraryTypeID     primitive.ObjectID `bson:"library_type_id" json:"library_type_id"`
	LibraryTemplateID primitive.ObjectID `bson:"library_template_id" json:"library_template_id"`
	LibraryType       *LibraryType       `bson:"-" json:"library_type"`
	LibraryTemplate   *LibraryTemplate   `bson:"-" json:"library_template"`
}

type LibraryTemplate struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name     string `bson:"name" json:"name"`
	Template string `bson:"template" json:"template"`
}

type LibraryType struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Name string `bson:"name" json:"name"`
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
	Overrides      *Overrides         `bson:"overrides" json:"overrides"`
	Paths          []*Path            `bson:"paths" json:"paths"`
	Cover          string             `bson:"-" json:"cover"`
	Background     string             `bson:"-" json:"background"`
	HasOverrides   bool               `bson:"-" json:"has_overrides"`
	Status         string             `bson:"status" json:"status"`
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
	ID           string `bson:"id" json:"id"`
	Name         string `bson:"name" json:"name"`
	Category     string `bson:"category" json:"category"`
	Dir          string `bson:"dir" json:"dir"`
	FinalDir     string `bson:"final_dir" json:"final_dir"`
	File         string `bson:"file" json:"file"`
	Status       string `bson:"status" json:"status"`
	StatusDetail string `bson:"status_detail" json:"status_detail"`
	StatusPar    string `bson:"status_par" json:"status_par"`
	StatusUnpack string `bson:"status_unpack" json:"status_unpack"`
}

type Overrides struct { // struct
	SeasonNumber   string `bson:"season_number" json:"season_number"`
	EpisodeNumber  string `bson:"episode_number" json:"episode_number"`
	AbsoluteNumber string `bson:"absolute_number" json:"absolute_number"`
}

type Path struct { // struct
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type       primitive.Symbol   `bson:"type" json:"type"`
	Tag        string             `bson:"tag" json:"tag"`
	Remote     string             `bson:"remote" json:"remote"`
	Local      string             `bson:"local" json:"local"`
	Extension  string             `bson:"extension" json:"extension"`
	Size       int64              `bson:"size" json:"size"`
	Resolution int                `bson:"resolution" json:"resolution"`
	Bitrate    int                `bson:"bitrate" json:"bitrate"`
	Checksum   string             `bson:"checksum" json:"checksum"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	Old        bool               `bson:"old" json:"old"`
	Rename     bool               `bson:"rename" json:"rename"`
}

type Pin struct { // model
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	//CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	//UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Pin        int    `bson:"pin" json:"id" grimoire:"index"`
	Code       string `bson:"code" json:"code" grimoire:"index"`
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
	Title    string `bson:"title" json:"title" grimoire:"index"`
	User     string `bson:"user" json:"user" grimoire:"index"`
	Type     string `bson:"type" json:"type"`
	Source   string `bson:"source" json:"source" grimoire:"index"`
	SourceID string `bson:"source_id" json:"source_id" grimoire:"index"`
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
	Name  string `bson:"name" json:"name" grimoire:"index"`
	Email string `bson:"email" json:"email" grimoire:"index"`
	Token string `bson:"token" json:"token"`
	Thumb string `bson:"thumb" json:"thumb"`
	Home  bool   `bson:"home" json:"home"`
	Admin bool   `bson:"admin" json:"admin"`
}

type Wanted struct { // struct
	Names    []string `bson:"names" json:"names"`
	Episodes []string `bson:"episodes" json:"episodes"`
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
