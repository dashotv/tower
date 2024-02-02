package plex

import (
	"crypto/tls"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	defaultMetaURL   = "https://metadata.provider.plex.tv"
	defaultPlexTVURL = "https://plex.tv/api/v2"
	applicationXml   = "application/xml"
	applicationJson  = "application/json"
)

const (
	PlexLibraryTypeUnknown = iota
	PlexLibraryTypeMovie
	PlexLibraryTypeShow
)

func New(opt *ClientOptions) *Client {
	c := &Client{
		URL:               opt.URL,
		Token:             opt.Token,
		MachineIdentifier: opt.MachineIdentifier,
		ClientIdentifier:  opt.ClientIdentifier,
		Product:           opt.Product,
		Device:            opt.Device,
		AppName:           opt.AppName,
		MetadataURL:       opt.MetadataURL,
		PlexTVURL:         opt.PlexTVURL,
	}
	if c.MetadataURL == "" {
		c.MetadataURL = defaultMetaURL
	}
	if c.PlexTVURL == "" {
		c.PlexTVURL = defaultPlexTVURL
	}

	c.Headers = map[string]string{
		"X-Plex-Token":             c.Token,
		"strong":                   "true",
		"Plex-Container-Size":      "50",
		"X-Plex-Container-Start":   "0",
		"X-Plex-Product":           c.AppName,
		"X-Plex-Client-Identifier": c.Device,
		"Accept":                   applicationJson,
		"ContentType":              applicationJson,
	}
	data := url.Values{}
	data.Set("strong", "true")
	data.Set("X-Plex-Client-Identifier", c.ClientIdentifier)
	data.Set("X-Plex-Product", c.Product)
	data.Set("X-Plex-Token", c.Token)
	c.data = data

	c.server = resty.New().SetBaseURL(c.URL).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	c.plextv = resty.New().SetBaseURL(c.PlexTVURL)
	c.metadata = resty.New().SetBaseURL(c.MetadataURL)

	return c
}

type ClientOptions struct {
	URL   string
	Token string

	MachineIdentifier string
	ClientIdentifier  string
	Product           string
	Device            string
	AppName           string
	MetadataURL       string
	PlexTVURL         string
}

type Client struct {
	URL   string
	Token string

	MachineIdentifier string
	ClientIdentifier  string
	Product           string
	Device            string
	AppName           string
	MetadataURL       string
	PlexTVURL         string
	Headers           map[string]string

	data     url.Values
	server   *resty.Client
	plextv   *resty.Client
	metadata *resty.Client
}

func (p *Client) _plextv() *resty.Request {
	return p.plextv.R().SetHeaders(p.Headers)
}
func (p *Client) _server() *resty.Request {
	return p.server.R().SetHeaders(p.Headers)
}
func (p *Client) _metadata() *resty.Request {
	return p.metadata.R().SetHeaders(p.Headers)
}

func (p *Client) GetUser(token string) (*PlexUser, error) {
	user := &PlexUser{}
	resp, err := p._plextv().SetResult(user).
		SetHeader("X-Plex-Token", token).Get("/user")
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("pin not authorized: %s", resp.Status())
	}

	return user, nil
}

type PlexLibraryMetadataContainer struct {
	MediaContainer struct {
		Size     int64                  `json:"size"`
		Metadata []*PlexLibraryMetadata `json:"Metadata"`
	} `json:"MediaContainer"`
}
type PlexLibraryMetadata struct {
	Key          string `json:"key"`
	RatingKey    string `json:"ratingKey"`
	Leaves       int    `json:"leafCount"`
	Viewed       int    `json:"viewedLeafCount"`
	LastViewedAt int64  `json:"lastViewedAt"`
}

func (p *Client) GetMetadataByKey(key string) (string, error) {
	resp, err := p._server().SetFormDataFromValues(p.data).Get("/library/metadata/" + key)
	if err != nil {
		return "", errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return "", errors.Errorf("failed to get metadata: %s", resp.Status())
	}
	return resp.String(), nil
}
func (p *Client) GetViewedByKey(key string) (*PlexLibraryMetadata, error) {
	m := &PlexLibraryMetadataContainer{}
	resp, err := p._server().SetResult(m).SetFormDataFromValues(p.data).Get("/library/metadata/" + key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get viewed: %s", resp.Status())
	}
	return m.MediaContainer.Metadata[0], nil
}
