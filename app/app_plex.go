package app

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

const (
	applicationXml  = "application/xml"
	applicationJson = "application/json"
)
const (
	PlexLibraryTypeUnknown = iota
	PlexLibraryTypeMovie
	PlexLibraryTypeShow
)

func init() {
	initializers = append(initializers, setupPlex)
}

func setupPlex(app *Application) error {
	plex := &Plex{
		URL: &PlexURLs{
			PlexTV:   app.Config.PlexTvURL,
			Server:   app.Config.PlexServerURL,
			Metadata: app.Config.PlexMetaURL,
		},
		Clients:    &PlexClients{},
		Identifier: app.Config.PlexClientIdentifier,
		Product:    app.Config.PlexAppName,
		Device:     app.Config.PlexDevice,
		Headers: map[string]string{
			"Plex-Container-Size":      "50",
			"X-Plex-Container-Start":   "0",
			"X-Plex-Product":           app.Config.PlexAppName,
			"X-Plex-Client-Identifier": app.Config.PlexDevice,
			"strong":                   "true",
			"Accept":                   applicationJson,
			"ContentType":              applicationJson,
			"X-Plex-Token":             app.Config.PlexToken,
		},
	}

	plex.Clients.PlexTV = resty.New().SetBaseURL(plex.URL.PlexTV)
	plex.Clients.Server = resty.New().SetBaseURL(plex.URL.Server)
	plex.Clients.Metadata = resty.New().SetBaseURL(plex.URL.Metadata)

	data := url.Values{}
	data.Set("strong", "true")
	data.Set("X-Plex-Client-Identifier", plex.Identifier)
	data.Set("X-Plex-Product", plex.Product)
	data.Set("X-Plex-Token", app.Config.PlexToken)
	plex.data = data

	app.Plex = plex
	return nil
}

type PlexURLs struct {
	PlexTV   string
	Server   string
	Metadata string
}
type PlexClients struct {
	PlexTV   *resty.Client
	Server   *resty.Client
	Metadata *resty.Client
}
type Plex struct {
	URL        *PlexURLs
	Clients    *PlexClients
	Identifier string
	Product    string
	Device     string
	Headers    map[string]string
	data       url.Values
}

func (p *Plex) plextv() *resty.Request {
	return p.Clients.PlexTV.R().SetHeaders(p.Headers)
}
func (p *Plex) server() *resty.Request {
	return p.Clients.Server.R().SetHeaders(p.Headers)
}
func (p *Plex) metadata() *resty.Request {
	return p.Clients.Metadata.R().SetHeaders(p.Headers)
}

// CreatePin returns a new pin from the plex api
func (p *Plex) CreatePin() (*Pin, error) {
	pin := &Pin{}
	resp, err := p.plextv().SetResult(pin).SetQueryParamsFromValues(p.data).Post("/pins")
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to create pin: %s", resp.Status())
	}
	app.Log.Debugf("create pin body: %s", resp.String())
	return pin, nil
}

func (p *Plex) CheckPin(pin *Pin) (bool, error) {
	params := url.Values{}
	params.Set("code", pin.Code)
	params.Set("X-Plex-Client-Identifier", app.Config.PlexClientIdentifier)

	newPin := &Pin{}
	resp, err := p.plextv().SetResult(newPin).
		SetHeader("code", pin.Code).
		SetQueryParamsFromValues(params).
		Get(fmt.Sprintf("/pins/%d", pin.Pin))
	if err != nil {
		return false, errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return false, errors.Errorf("pin not authorized: %s", resp.Status())
	}
	if newPin.Token == "" {
		return false, errors.Errorf("pin not authorized: token is empty")
	}

	pin.Token = newPin.Token
	pin.Product = newPin.Product
	pin.Identifier = newPin.Identifier

	err = app.DB.Pin.Update(pin)
	if err != nil {
		return false, errors.Wrap(err, "failed to update token")
	}

	return true, nil
}

func (p *Plex) GetUser(token string) (*PlexUser, error) {
	user := &PlexUser{}
	resp, err := p.plextv().SetResult(user).
		SetHeader("X-Plex-Token", token).Get("/user")
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("pin not authorized: %s", resp.Status())
	}

	return user, nil
}

func (p *Plex) getAuthUrl(pin *Pin) string {
	base := "https://app.plex.tv/auth/#?"
	data := url.Values{}
	data.Set("clientID", p.Identifier)
	data.Set("code", pin.Code)
	data.Set("forwardUrl", fmt.Sprintf("%s/auth?pin=%d", app.Config.Plex, pin.Pin))
	data.Set("context[device][product]", p.Product)
	data.Set("context[device][version]", "0.1.0")
	data.Set("context[device][deviceName]", p.Device)
	return base + data.Encode()
}

/*
Plex Pin response:
{
	"id": 00000000000,
	"code": "adladoqienbfquboqoqiobeoi",
	"product": "DashoTV",
	"trusted": false,
	"qr": "https://plex.tv/api/v2/pins/qr/adladoqienbfquboqoqiobeoi",
	"clientIdentifier": "dashotv-web",
	"location": {
		"code": "US",
		"european_union_member": false,
		"continent_code": "NA",
		"country": "United States",
		"city": "San Francisco",
		"time_zone": "America/Los_Angeles",
		"postal_code": "94124",
		"in_privacy_restricted_country": false,
		"subdivisions": "California",
		"coordinates": "37.7308, -122.3838"
	},
	"expiresIn": 1800,
	"createdAt": "2023-10-14T23:53:35Z",
	"expiresAt": "2023-10-15T00:23:35Z",
	"authToken": null, // set after auth
	"newRegistration": null // set after auth
}
*/
/*
{
        "allowSync": true,
        "art": "/:/resources/show-fanart.jpg",
        "composite": "/library/sections/2/composite/1704818054",
        "filters": true,
        "refreshing": true,
        "thumb": "/:/resources/show.png",
        "key": "2",
        "type": "show",
        "title": "TV Shows",
        "agent": "tv.plex.agents.series",
        "scanner": "Plex TV Series",
        "language": "en-US",
        "uuid": "e35a54e6-79ff-4cde-98c6-c05c1a09c821",
        "updatedAt": 1704818112,
        "createdAt": 1384056048,
        "scannedAt": 1704818054,
        "content": true,
        "directory": true,
        "contentChangedAt": 41632394,
        "hidden": 0,
        "Location": [
          {
            "id": 29,
            "path": "/mnt/media/tv"
          }
        ]
      },
*/
type PlexLibraryLocation struct {
	ID   int64  `json:"id"`
	Path string `json:"path"`
}

type PlexLibrary struct {
	AllowSync        bool                   `json:"allowSync"`
	Art              string                 `json:"art"`
	Composite        string                 `json:"composite"`
	Filters          bool                   `json:"filters"`
	Refreshing       bool                   `json:"refreshing"`
	Thumb            string                 `json:"thumb"`
	Key              string                 `json:"key"`
	Type             string                 `json:"type"`
	Title            string                 `json:"title"`
	Agent            string                 `json:"agent"`
	Scanner          string                 `json:"scanner"`
	Language         string                 `json:"language"`
	UUID             string                 `json:"uuid"`
	UpdatedAt        int64                  `json:"updatedAt"`
	CreatedAt        int64                  `json:"createdAt"`
	ScannedAt        int64                  `json:"scannedAt"`
	Content          bool                   `json:"content"`
	Directory        bool                   `json:"directory"`
	ContentChangedAt int64                  `json:"contentChangedAt"`
	Hidden           int64                  `json:"hidden"`
	Locations        []*PlexLibraryLocation `json:"Location"`
}

type PlexLibraries struct {
	MediaContainer struct {
		Size        int64          `json:"size"`
		Directories []*PlexLibrary `json:"Directory"`
	} `json:"MediaContainer"`
}

func (p *Plex) GetLibraries() ([]*PlexLibrary, error) {
	dest := &PlexLibraries{}
	resp, err := p.server().SetResult(dest).SetFormDataFromValues(p.data).Get("/library/sections")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get libraries: %s", resp.Status())
	}

	return dest.MediaContainer.Directories, nil
}

func (p *Plex) LibraryType(section string) (int, error) {
	t, err := p.LibraryTypeName(section)
	if err != nil {
		return 0, err
	}
	if t == "" {
		return 0, errors.Errorf("library section %s not found", section)
	}

	id := p.LibraryTypeID(t)
	if id == PlexLibraryTypeUnknown {
		return 0, errors.Errorf("library section %s unknown", section)
	}

	return id, nil
}

func (p *Plex) LibraryTypeName(section string) (string, error) {
	resp, err := p.GetLibraries()
	if err != nil {
		return "", err
	}
	for _, r := range resp {
		if r.Key == section {
			return r.Type, nil
		}
	}
	return "", errors.Errorf("library section %s not found", section)
}

func (p *Plex) LibraryTypeID(t string) int {
	switch t {
	case "movie":
		return PlexLibraryTypeMovie
	case "show":
		return PlexLibraryTypeShow
	default:
		return PlexLibraryTypeUnknown
	}
}

func (p *Plex) Search(query, section string) ([]SearchMetadata, error) {
	id, err := p.LibraryType(section)
	if err != nil {
		return nil, err
	}

	dest := &PlexSearch{}
	path := fmt.Sprintf("/library/sections/%s/search", section)

	params := url.Values{}
	params.Set("X-Plex-Token", app.Config.PlexToken)
	params.Set("title", query)
	params.Set("type", fmt.Sprintf("%d", id))
	params.Set("limit", "25")
	params.Set("sort", "createdAt:desc")

	resp, err := p.server().SetResult(dest).
		SetQueryParamsFromValues(params).
		Get(path)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get search: %s", resp.Status())
	}

	// app.Log.Debugf("search req url: %s", resp.Request.URL)
	// app.Log.Debugf("search result: %s", resp.String())
	return dest.MediaContainer.Metadata, nil
}

type PlexSearch struct {
	MediaContainer struct {
		Size         int64            `json:"size"`
		SectionID    int64            `json:"sectionID"`
		AllowSync    bool             `json:"allowSync"`
		Art          string           `json:"art"`
		Identifier   string           `json:"identifier"`
		SectionTitle string           `json:"librarySectionTitle"`
		SectionUUID  string           `json:"librarySectionUUID"`
		Title        string           `json:"title1"`
		Subtitle     string           `json:"title2"`
		Metadata     []SearchMetadata `json:"Metadata,omitempty"`
	} `json:"MediaContainer"`
}

type SearchMetadata struct {
	LibrarySectionTitle   string     `json:"librarySectionTitle"`
	Score                 string     `json:"score"`
	RatingKey             string     `json:"ratingKey"`
	Key                   string     `json:"key"`
	GUID                  string     `json:"guid"`
	Slug                  *string    `json:"slug,omitempty"`
	Studio                *string    `json:"studio,omitempty"`
	Type                  string     `json:"type"`
	Title                 string     `json:"title"`
	LibrarySectionID      int64      `json:"librarySectionID"`
	LibrarySectionKey     string     `json:"librarySectionKey"`
	ContentRating         string     `json:"contentRating"`
	Summary               string     `json:"summary"`
	Index                 *int64     `json:"index,omitempty"`
	AudienceRating        *float64   `json:"audienceRating,omitempty"`
	ViewCount             *int64     `json:"viewCount,omitempty"`
	LastViewedAt          *int64     `json:"lastViewedAt,omitempty"`
	Year                  *int64     `json:"year,omitempty"`
	Tagline               *string    `json:"tagline,omitempty"`
	Thumb                 string     `json:"thumb"`
	Art                   string     `json:"art"`
	Theme                 *string    `json:"theme,omitempty"`
	Duration              int64      `json:"duration"`
	OriginallyAvailableAt string     `json:"originallyAvailableAt"`
	LeafCount             *int64     `json:"leafCount,omitempty"`
	ViewedLeafCount       *int64     `json:"viewedLeafCount,omitempty"`
	ChildCount            *int64     `json:"childCount,omitempty"`
	AddedAt               int64      `json:"addedAt"`
	UpdatedAt             int64      `json:"updatedAt"`
	AudienceRatingImage   *string    `json:"audienceRatingImage,omitempty"`
	Genre                 []Country  `json:"Genre,omitempty"`
	Country               []Country  `json:"Country,omitempty"`
	Role                  []Country  `json:"Role,omitempty"`
	Location              []Location `json:"Location,omitempty"`
	SkipCount             *int64     `json:"skipCount,omitempty"`
	SeasonCount           *int64     `json:"seasonCount,omitempty"`
	Field                 []Field    `json:"Field,omitempty"`
	Rating                *float64   `json:"rating,omitempty"`
	PrimaryExtraKey       *string    `json:"primaryExtraKey,omitempty"`
	RatingImage           *string    `json:"ratingImage,omitempty"`
	Media                 []Media    `json:"Media,omitempty"`
	Director              []Country  `json:"Director,omitempty"`
	Writer                []Country  `json:"Writer,omitempty"`
	ChapterSource         *string    `json:"chapterSource,omitempty"`
	ParentRatingKey       *string    `json:"parentRatingKey,omitempty"`
	GrandparentRatingKey  *string    `json:"grandparentRatingKey,omitempty"`
	ParentGUID            *string    `json:"parentGuid,omitempty"`
	GrandparentGUID       *string    `json:"grandparentGuid,omitempty"`
	TitleSort             *string    `json:"titleSort,omitempty"`
	GrandparentKey        *string    `json:"grandparentKey,omitempty"`
	ParentKey             *string    `json:"parentKey,omitempty"`
	GrandparentTitle      *string    `json:"grandparentTitle,omitempty"`
	ParentTitle           *string    `json:"parentTitle,omitempty"`
	OriginalTitle         *string    `json:"originalTitle,omitempty"`
	ParentIndex           *int64     `json:"parentIndex,omitempty"`
	ParentYear            *int64     `json:"parentYear,omitempty"`
	ParentThumb           *string    `json:"parentThumb,omitempty"`
	GrandparentThumb      *string    `json:"grandparentThumb,omitempty"`
	GrandparentArt        *string    `json:"grandparentArt,omitempty"`
	GrandparentTheme      *string    `json:"grandparentTheme,omitempty"`
	GrandparentSlug       *string    `json:"grandparentSlug,omitempty"`
}

type Country struct {
	Tag string `json:"tag"`
}

type Field struct {
	Locked bool   `json:"locked"`
	Name   string `json:"name"`
}

type Location struct {
	Path string `json:"path"`
}

type Media struct {
	ID                    int64   `json:"id"`
	Duration              int64   `json:"duration"`
	Bitrate               int64   `json:"bitrate"`
	Width                 int64   `json:"width"`
	Height                int64   `json:"height"`
	AspectRatio           float64 `json:"aspectRatio"`
	AudioChannels         int64   `json:"audioChannels"`
	AudioCodec            string  `json:"audioCodec"`
	VideoCodec            string  `json:"videoCodec"`
	VideoResolution       string  `json:"videoResolution"`
	Container             string  `json:"container"`
	VideoFrameRate        string  `json:"videoFrameRate"`
	AudioProfile          *string `json:"audioProfile,omitempty"`
	VideoProfile          string  `json:"videoProfile"`
	Part                  []Part  `json:"Part"`
	OptimizedForStreaming *int64  `json:"optimizedForStreaming,omitempty"`
	Has64BitOffsets       *bool   `json:"has64bitOffsets,omitempty"`
}

type Part struct {
	ID                    int64   `json:"id"`
	Key                   string  `json:"key"`
	Duration              int64   `json:"duration"`
	File                  string  `json:"file"`
	Size                  int64   `json:"size"`
	AudioProfile          *string `json:"audioProfile,omitempty"`
	Container             string  `json:"container"`
	VideoProfile          string  `json:"videoProfile"`
	Has64BitOffsets       *bool   `json:"has64bitOffsets,omitempty"`
	OptimizedForStreaming *bool   `json:"optimizedForStreaming,omitempty"`
}

type Style string

const (
	Shelf Style = "shelf"
)

//
// type Hub struct {
// 	Title         string        `json:"title"`
// 	Type          string        `json:"type"`
// 	HubIdentifier string        `json:"hubIdentifier"`
// 	Context       string        `json:"context"`
// 	Size          int64         `json:"size"`
// 	More          bool          `json:"more"`
// 	Style         Style         `json:"style"`
// 	Metadata      []HubMetadata `json:"Metadata,omitempty"`
// 	Directory     []Directory   `json:"Directory,omitempty"`
// }
//
// type Directory struct {
// 	Key                 string `json:"key"`
// 	LibrarySectionID    int64  `json:"librarySectionID"`
// 	LibrarySectionKey   string `json:"librarySectionKey"`
// 	LibrarySectionTitle string `json:"librarySectionTitle"`
// 	LibrarySectionType  int64  `json:"librarySectionType"`
// 	Reason              string `json:"reason"`
// 	ReasonID            int64  `json:"reasonID"`
// 	ReasonTitle         string `json:"reasonTitle"`
// 	Score               string `json:"score"`
// 	Type                string `json:"type"`
// 	ID                  int64  `json:"id"`
// 	Filter              string `json:"filter"`
// 	Tag                 string `json:"tag"`
// 	TagType             int64  `json:"tagType"`
// 	Thumb               string `json:"thumb"`
// 	Art                 string `json:"art"`
// 	Count               int64  `json:"count"`
// 	GUID                string `json:"guid"`
// 	Summary             string `json:"summary"`
// }

type PlexCollectionCreate struct {
	MediaContainer struct {
		Directory []struct {
			RatingKey string `json:"ratingKey"`
			Key       string `json:"key"`
			Guid      string `json:"guid"`
			Type      string `json:"type"`
			Title     string `json:"title"`
			Subtype   string `json:"subtype"`
			Summary   string `json:"summary"`
			Thumb     string `json:"thumb"`
			AddedAt   int64  `json:"addedAt"`
			UpdatedAt int64  `json:"updatedAt"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}

func (p *Plex) CreateCollection(title, section, firstKey string) (*PlexCollectionCreate, error) {
	id, err := p.LibraryType(section)
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("X-Plex-Token", app.Config.PlexToken)
	data.Set("title", title)
	data.Set("sectionId", section)
	data.Set("type", fmt.Sprintf("%d", id))
	data.Set("smart", "0")
	data.Set("uri", fmt.Sprintf("server://%s/com.plexapp.plugins.library/library/metadata/%s", app.Config.PlexMachineIdentifier, firstKey))

	dest := &PlexCollectionCreate{}
	resp, err := p.server().
		SetResult(dest).
		SetQueryParamsFromValues(data).
		Post("/library/collections")
	if err != nil {
		return nil, err
	}
	app.Log.Debugf("create collection req url: %s", resp.Request.URL)
	app.Log.Debugf("create collection response: %s", resp.String())
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to create collection: %s", resp.Status())
	}

	return dest, nil
}

type PlexLibrariesCollectionResponse struct {
	MediaContainer struct {
		Size         int64             `json:"size"`
		AllowSync    bool              `json:"allowSync"`
		Identifier   string            `json:"identifier"`
		LibraryID    int64             `json:"librarySectionID"`
		LibraryTitle string            `json:"librarySectionTitle"`
		LibraryUUID  string            `json:"librarySectionUUID"`
		Title        string            `json:"title1"`
		Subtitle     string            `json:"title2"`
		Metadata     []*PlexCollection `json:"Metadata,omitempty"`
	} `json:"MediaContainer"`
}

type PlexCollectionResponse struct {
	MediaContainer struct {
		Size         int64             `json:"size"`
		AllowSync    bool              `json:"allowSync"`
		Identifier   string            `json:"identifier"`
		LibraryID    int64             `json:"librarySectionID"`
		LibraryTitle string            `json:"librarySectionTitle"`
		LibraryUUID  string            `json:"librarySectionUUID"`
		Title        string            `json:"title1"`
		Subtitle     string            `json:"title2"`
		Directory    []*PlexCollection `json:"Metadata,omitempty"`
	} `json:"MediaContainer"`
}
type PlexCollection struct {
	RatingKey    string                 `json:"ratingKey"`
	Key          string                 `json:"key"`
	GUID         string                 `json:"guid"`
	Type         string                 `json:"type"`
	Title        string                 `json:"title"`
	LibraryID    int64                  `json:"librarySectionID"`
	LibraryTitle string                 `json:"librarySectionTitle"`
	LibraryKey   string                 `json:"librarySectionKey"`
	Subtype      string                 `json:"subtype"`
	Summary      string                 `json:"summary"`
	Thumb        string                 `json:"thumb"`
	AddedAt      int64                  `json:"addedAt"`
	UpdatedAt    int64                  `json:"updatedAt"`
	ChildCount   string                 `json:"childCount"`
	MaxYear      string                 `json:"maxYear"`
	MinYear      string                 `json:"minYear"`
	Children     []*PlexCollectionChild `json:"children,omitempty"`
}
type PlexCollectionChildrenResponse struct {
	MediaContainer struct {
		Size      int64                  `json:"size"`
		Directory []*PlexCollectionChild `json:"Metadata,omitempty"`
	} `json:"MediaContainer"`
}
type PlexCollectionChild struct {
	RatingKey    string `json:"ratingKey"`
	Key          string `json:"key"`
	GUID         string `json:"guid"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	LibraryID    int64  `json:"librarySectionID"`
	LibraryTitle string `json:"librarySectionTitle"`
	LibraryKey   string `json:"librarySectionKey"`
	Summary      string `json:"summary"`
	Thumb        string `json:"thumb"`
	AddedAt      int64  `json:"addedAt"`
	UpdatedAt    int64  `json:"updatedAt"`
}

func (p *Plex) DeleteCollection(ratingKey string) error {
	data := url.Values{}
	data.Set("X-Plex-Token", app.Config.PlexToken)

	resp, err := p.server().
		SetQueryParamsFromValues(data).
		Delete(fmt.Sprintf("/library/collections/%s", ratingKey))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return errors.Errorf("failed to delete collection: %s", resp.Status())
	}

	return nil
}

func (p *Plex) ListCollections(section string) ([]*PlexCollection, error) {
	data := url.Values{}
	data.Set("X-Plex-Token", app.Config.PlexToken)

	dest := &PlexLibrariesCollectionResponse{}
	resp, err := p.server().
		SetResult(dest).
		SetQueryParamsFromValues(p.data).
		Get("/library/sections/" + section + "/collections")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get collections: %s", resp.Status())
	}

	return dest.MediaContainer.Metadata, nil
}

func (p *Plex) GetCollection(ratingKey string) (*PlexCollection, error) {
	data := url.Values{}
	data.Set("X-Plex-Token", app.Config.PlexToken)

	dest := &PlexCollectionResponse{}
	resp, err := p.server().
		SetResult(dest).
		SetQueryParamsFromValues(p.data).
		Get("/library/collections/" + ratingKey)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to update collection: %s", resp.Status())
	}
	if len(dest.MediaContainer.Directory) != 1 {
		return nil, errors.Errorf("api response found %d directories, wanted 1", len(dest.MediaContainer.Directory))
	}

	children, err := p.GetCollectionChildren(ratingKey)
	if err != nil {
		return nil, err
	}

	r := dest.MediaContainer.Directory[0]
	r.Children = children

	return r, nil
}

func (p *Plex) GetCollectionChildren(ratingKey string) ([]*PlexCollectionChild, error) {
	data := url.Values{}
	data.Set("X-Plex-Token", app.Config.PlexToken)

	dest := &PlexCollectionChildrenResponse{}
	resp, err := p.server().
		SetResult(dest).
		SetQueryParamsFromValues(p.data).
		Get("/library/collections/" + ratingKey + "/children")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get collection children: %s", resp.Status())
	}

	// app.Log.Debugf("collection children: %s", resp.String())

	return dest.MediaContainer.Directory, nil
}

func (p *Plex) UpdateCollection(section, ratingKey string, keys []string) error {
	existing, err := p.GetCollection(ratingKey)
	if err != nil {
		return err
	}

	existingKeys := lo.Map(existing.Children, func(c *PlexCollectionChild, i int) string {
		return c.RatingKey
	})

	add, remove := lo.Difference(keys, existingKeys)
	if len(add) > 0 {
		app.Log.Debugf("adding %d items to collection: %+v", len(add), add)
		for _, k := range add {
			if err := p.addCollectionItem(ratingKey, k); err != nil {
				return err
			}
		}
	}
	if len(remove) > 0 {
		app.Log.Debugf("removing %d items from collection: %+v", len(remove), remove)
		for _, k := range remove {
			if err := p.removeCollectionItem(ratingKey, k); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Plex) addCollectionItem(ratingKey, newKey string) error {
	data := url.Values{}
	data.Set("X-Plex-Token", app.Config.PlexToken)
	data.Set("uri", fmt.Sprintf("server://%s/com.plexapp.plugins.library/library/metadata/%s", app.Config.PlexMachineIdentifier, newKey))

	resp, err := p.server().
		SetQueryParamsFromValues(data).
		Put("/library/collections/" + ratingKey + "/items")
	app.Log.Debugf("addCollectionItem req url: %s", resp.Request.URL)
	app.Log.Debugf("addCollectionItem result: %s", resp.String())
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return errors.Errorf("failed to add to collection: %s", resp.Status())
	}

	return nil
}

func (p *Plex) removeCollectionItem(ratingKey, rmKey string) error {
	data := url.Values{}
	data.Set("X-Plex-Token", app.Config.PlexToken)
	data.Set("excludeAllLeaves", "1")

	resp, err := p.server().
		SetQueryParamsFromValues(data).
		Delete(fmt.Sprintf("/library/collections/%s/children/%s", ratingKey, rmKey))
	app.Log.Debugf("removeCollectionItem req url: %s", resp.Request.URL)
	app.Log.Debugf("removeCollectionItem result: %s", resp.String())
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return errors.Errorf("failed to remove from collection: %s", resp.Status())
	}

	return nil
}

// func (p *Plex) GetCollections(token string) (*PlexCollections, error) {
// 	dest := &PlexCollections{}
// 	resp, err := p.metadata().SetResult(dest).SetHeader("X-Plex-Token", token).Get("/library/sections")
// 	if err != nil {
// 		return dest, err
// 	}
// 	if !resp.IsSuccess() {
// 		return dest, errors.Errorf("failed to get collections: %s", resp.Status())
// 	}
//
// 	return dest, nil
// }

type WatchlistOpts struct {
	Filter string // all, or ???
	Sort   string // ???
	Type   string // library type? movie, show, episode, artist, album, track?
}

func (p *Plex) GetWatchlist(token string) (*PlexWatchlist, error) {
	dest := &PlexWatchlist{}
	opts := &WatchlistOpts{Filter: "all"}
	u := fmt.Sprintf("/library/sections/watchlist/%s", opts.Filter)
	resp, err := p.metadata().SetResult(dest).SetHeader("X-Plex-Token", token).Get(u)
	if err != nil {
		return dest, err
	}
	if !resp.IsSuccess() {
		return dest, errors.Errorf("failed to get watchlist: %s", resp.Status())
	}

	return dest, nil
}

func (p *Plex) GetWatchlistDetail(token string, w *PlexWatchlist) ([]*PlexWatchlistDetail, error) {
	out := []*PlexWatchlistDetail{}
	for _, d := range w.MediaContainer.Metadata {
		dest := &PlexWatchlistDetail{}
		resp, err := p.metadata().
			SetResult(dest).
			SetHeader("X-Plex-Token", token).
			Get("/library/metadata/" + d.RatingKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to make watchlistdetails request")
		}
		if !resp.IsSuccess() {
			return nil, errors.Errorf("failed to get watchlistdetails: %s", resp.Status())
		}
		out = append(out, dest)
	}
	return out, nil
}

type PlexUser struct {
	ID       int64  `json:"id"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	// Title             string      `json:"title"`
	Email string `json:"email"`
	// FriendlyName      string      `json:"friendlyName"`
	// Locale            interface{} `json:"locale"`
	Confirmed bool  `json:"confirmed"`
	JoinedAt  int64 `json:"joinedAt"`
	// EmailOnlyAuth     bool        `json:"emailOnlyAuth"`
	// HasPassword       bool        `json:"hasPassword"`
	// Protected         bool        `json:"protected"`
	Thumb string `json:"thumb"`
	// AuthToken         string      `json:"authToken"`
	// MailingListStatus string      `json:"mailingListStatus"`
	// MailingListActive bool        `json:"mailingListActive"`
	// ScrobbleTypes     string      `json:"scrobbleTypes"`
	// Country           string      `json:"country"`
	// Pin               string      `json:"pin"`
	// Subscription      struct {
	// 	Active         bool     `json:"active"`
	// 	SubscribedAt   string   `json:"subscribedAt"`
	// 	Status         string   `json:"status"`
	// 	PaymentService string   `json:"paymentService"`
	// 	Plan           string   `json:"plan"`
	// 	Features       []string `json:"features"`
	// } `json:"subscription"`
	// SubscriptionDescription string `json:"subscriptionDescription"`
	// Restricted              bool   `json:"restricted"`
	// Anonymous               bool   `json:"anonymous"`
	Home bool `json:"home"`
	// Guest                   bool   `json:"guest"`
	HomeSize  int64 `json:"homeSize"`
	HomeAdmin bool  `json:"homeAdmin"`
	// MaxHomeSize             int64  `json:"maxHomeSize"`
	// RememberExpiresAt       int64  `json:"rememberExpiresAt"`
	// Profile                 struct {
	// 	AutoSelectAudio              bool   `json:"autoSelectAudio"`
	// 	DefaultAudioLanguage         string `json:"defaultAudioLanguage"`
	// 	DefaultSubtitleLanguage      string `json:"defaultSubtitleLanguage"`
	// 	AutoSelectSubtitle           int64  `json:"autoSelectSubtitle"`
	// 	DefaultSubtitleAccessibility int64  `json:"defaultSubtitleAccessibility"`
	// 	DefaultSubtitleForced        int64  `json:"defaultSubtitleForced"`
	// } `json:"profile"`
	// Entitlements []string `json:"entitlements"`
	// Roles        []string `json:"roles"`
	// Services     []struct {
	// 	Identifier string  `json:"identifier"`
	// 	Endpoint   string  `json:"endpoint"`
	// 	Token      *string `json:"token"`
	// 	Secret     *string `json:"secret"`
	// 	Status     string  `json:"status"`
	// } `json:"services"`
	// AdsConsent           interface{} `json:"adsConsent"`
	// AdsConsentSetAt      interface{} `json:"adsConsentSetAt"`
	// AdsConsentReminderAt interface{} `json:"adsConsentReminderAt"`
	// ExperimentalFeatures bool        `json:"experimentalFeatures"`
	// TwoFactorEnabled     bool        `json:"twoFactorEnabled"`
	// BackupCodesCreated   bool        `json:"backupCodesCreated"`
}

type PlexWatchlist struct {
	MediaContainer struct {
		LibrarySectionID    string `json:"librarySectionID"`
		LibrarySectionTitle string `json:"librarySectionTitle"`
		Offset              int64  `json:"offset"`
		TotalSize           int64  `json:"totalSize"`
		Identifier          string `json:"identifier"`
		Size                int64  `json:"size"`
		Metadata            []struct {
			Art                   string   `json:"art"`
			Banner                string   `json:"banner"`
			GUID                  string   `json:"guid"`
			Key                   string   `json:"key"`
			Rating                float64  `json:"rating"`
			RatingKey             string   `json:"ratingKey"`
			Studio                string   `json:"studio"`
			Type                  string   `json:"type"`
			Theme                 *string  `json:"theme,omitempty"`
			Thumb                 string   `json:"thumb"`
			AddedAt               int64    `json:"addedAt"`
			Duration              int64    `json:"duration"`
			PublicPagesURL        string   `json:"publicPagesURL"`
			Slug                  string   `json:"slug"`
			UserState             bool     `json:"userState"`
			Title                 string   `json:"title"`
			OriginalTitle         *string  `json:"originalTitle,omitempty"`
			LeafCount             *int64   `json:"leafCount,omitempty"`
			ChildCount            *int64   `json:"childCount,omitempty"`
			ContentRating         string   `json:"contentRating"`
			OriginallyAvailableAt string   `json:"originallyAvailableAt"`
			Year                  int64    `json:"year"`
			RatingImage           string   `json:"ratingImage"`
			ImdbRatingCount       int64    `json:"imdbRatingCount"`
			Image                 []Image  `json:"Image"`
			PrimaryExtraKey       *string  `json:"primaryExtraKey,omitempty"`
			IsContinuingSeries    *bool    `json:"isContinuingSeries,omitempty"`
			Tagline               *string  `json:"tagline,omitempty"`
			AudienceRating        *float64 `json:"audienceRating,omitempty"`
			AudienceRatingImage   *string  `json:"audienceRatingImage,omitempty"`
			Subtype               *string  `json:"subtype,omitempty"`
			SkipChildren          *bool    `json:"skipChildren,omitempty"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}

type PlexWatchlistDetail struct {
	MediaContainer struct {
		Offset     int64  `json:"offset"`
		TotalSize  int64  `json:"totalSize"`
		Identifier string `json:"identifier"`
		Size       int64  `json:"size"`
		Metadata   []struct {
			Art                   string     `json:"art"`
			Banner                string     `json:"banner"`
			MetadatumGUID         string     `json:"guid"`
			Key                   string     `json:"key"`
			MetadatumRating       float64    `json:"rating"`
			RatingKey             string     `json:"ratingKey"`
			MetadatumStudio       string     `json:"studio"`
			Summary               string     `json:"summary"`
			Tagline               string     `json:"tagline"`
			Type                  string     `json:"type"`
			Theme                 string     `json:"theme"`
			Thumb                 string     `json:"thumb"`
			AddedAt               int64      `json:"addedAt"`
			Duration              int64      `json:"duration"`
			PublicPagesURL        string     `json:"publicPagesURL"`
			Slug                  string     `json:"slug"`
			UserState             bool       `json:"userState"`
			Title                 string     `json:"title"`
			LeafCount             int64      `json:"leafCount"`
			ChildCount            int64      `json:"childCount"`
			ContentRating         string     `json:"contentRating"`
			OriginallyAvailableAt string     `json:"originallyAvailableAt"`
			Year                  int64      `json:"year"`
			RatingImage           string     `json:"ratingImage"`
			ImdbRatingCount       int64      `json:"imdbRatingCount"`
			Image                 []Image    `json:"Image"`
			Genre                 []Genre    `json:"Genre"`
			GUID                  []GUID     `json:"Guid"`
			Country               []Country  `json:"Country"`
			Role                  []Role     `json:"Role"`
			Director              []Director `json:"Director"`
			Producer              []Director `json:"Producer"`
			Writer                []Director `json:"Writer"`
			Network               []Country  `json:"Network"`
			Rating                []Rating   `json:"Rating"`
			Similar               []Similar  `json:"Similar"`
			Studio                []Country  `json:"Studio"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}

type Director struct {
	ID        string  `json:"id"`
	Slug      string  `json:"slug"`
	Tag       string  `json:"tag"`
	Role      string  `json:"role"`
	Directory bool    `json:"directory"`
	Thumb     *string `json:"thumb,omitempty"`
}

type GUID struct {
	ID string `json:"id"`
}

type Genre struct {
	Filter      string  `json:"filter"`
	ID          string  `json:"id"`
	RatingKey   string  `json:"ratingKey"`
	Slug        string  `json:"slug"`
	Tag         string  `json:"tag"`
	Directory   bool    `json:"directory"`
	Context     string  `json:"context"`
	OriginalTag *string `json:"originalTag,omitempty"`
}

type Image struct {
	Alt  string `json:"alt"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Rating struct {
	Image string  `json:"image"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type Role struct {
	ID        string  `json:"id"`
	Order     int64   `json:"order"`
	Slug      string  `json:"slug"`
	Tag       string  `json:"tag"`
	Thumb     *string `json:"thumb,omitempty"`
	Role      *string `json:"role,omitempty"`
	Directory bool    `json:"directory"`
}

type Similar struct {
	GUID string `json:"guid"`
	Tag  string `json:"tag"`
}
