package app

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	applicationXml  = "application/xml"
	applicationJson = "application/json"
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

//	func (p *Plex) server() *resty.Request {
//		return p.Clients.Server.R().SetHeaders(p.Headers)
//	}
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
	newPin := &Pin{}
	resp, err := p.plextv().SetResult(newPin).
		SetHeader("code", pin.Code).
		SetQueryParamsFromValues(url.Values{"code": {pin.Code}}).
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

type Country struct {
	Tag string `json:"tag"`
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
