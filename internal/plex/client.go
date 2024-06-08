package plex

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"

	"github.com/dashotv/fae"
)

const (
	defaultMetaURL  = "https://metadata.provider.plex.tv"
	defaultTVURL    = "https://plex.tv/api"
	applicationXml  = "application/xml"
	applicationJson = "application/json"
)

const (
	LibraryTypeUnknown = iota
	LibraryTypeMovie
	LibraryTypeShow
)

func New(opt *ClientOptions) *Client {
	c := &Client{
		URL:               opt.URL,
		Token:             opt.Token,
		Debug:             opt.Debug,
		MachineIdentifier: opt.MachineIdentifier,
		ClientIdentifier:  opt.ClientIdentifier,
		Product:           opt.Product,
		Device:            opt.Device,
		AppName:           opt.AppName,
		MetadataURL:       opt.MetadataURL,
		TVURL:             opt.TVURL,
	}
	if c.MetadataURL == "" {
		c.MetadataURL = defaultMetaURL
	}
	if c.TVURL == "" {
		c.TVURL = defaultTVURL
	}

	c.Headers = map[string]string{
		"X-Plex-Token":             c.Token,
		"X-Plex-Container-Size":    "50",
		"X-Plex-Container-Start":   "0",
		"X-Plex-Product":           c.AppName,
		"X-Plex-Client-Identifier": c.Device,
		"strong":                   "true",
		"Accept":                   applicationJson,
		"Content-type":             applicationJson,
	}

	data := url.Values{}
	data.Set("strong", "true")
	data.Set("X-Plex-Client-Identifier", c.ClientIdentifier)
	data.Set("X-Plex-Product", c.Product)
	data.Set("X-Plex-Token", c.Token)
	c.data = data

	c.server = resty.New().SetDebug(c.Debug).SetBaseURL(c.URL).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	c.plextv = resty.New().SetDebug(c.Debug).SetBaseURL(c.TVURL + "/v2")
	c.plextvv1 = resty.New().SetDebug(c.Debug).SetBaseURL(c.TVURL)
	c.metadata = resty.New().SetDebug(c.Debug).SetBaseURL(c.MetadataURL)

	return c
}

type ClientOptions struct {
	URL   string
	Token string
	Debug bool

	MachineIdentifier string
	ClientIdentifier  string
	Product           string
	Device            string
	AppName           string
	MetadataURL       string
	TVURL             string
}

type Client struct {
	URL   string
	Token string
	Debug bool

	MachineIdentifier string
	ClientIdentifier  string
	Product           string
	Device            string
	AppName           string
	MetadataURL       string
	TVURL             string
	Headers           map[string]string

	data     url.Values
	server   *resty.Client
	plextv   *resty.Client
	plextvv1 *resty.Client
	metadata *resty.Client
	accounts []*User // cache of users
}

func (p *Client) _plextv() *resty.Request {
	return p.plextv.R().SetHeaders(p.Headers)
}
func (p *Client) _plextvv1() *resty.Request {
	return p.plextvv1.R().SetHeaders(p.Headers)
}
func (p *Client) _server() *resty.Request {
	return p.server.R().SetHeaders(p.Headers)
}
func (p *Client) _metadata() *resty.Request {
	return p.metadata.R().SetHeaders(p.Headers)
}

type Accounts struct {
	XMLName xml.Name `json:"MediaContainer" xml:"MediaContainer"`
	Size    int64    `json:"size" xml:"size,attr"`
	Total   int64    `json:"totalSize"  xml:"totalSize,attr"`
	Users   []*User  `json:"User" xml:"User"`
}

func (p *Client) GetAccounts() ([]*User, error) {
	if p.accounts != nil && len(p.accounts) > 0 {
		return p.accounts, nil
	}
	return p.GetAccountsUpdate()
}

func (p *Client) GetAccountsUpdate() ([]*User, error) {
	r := &Accounts{}
	resp, err := p._plextvv1().
		SetResult(r).
		Get("/users/")
	if err != nil {
		return nil, fae.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, fae.Errorf("getting accounts: %s", resp.Status())
	}

	p.accounts = r.Users
	return p.accounts, nil
}

func (p *Client) GetAccount(id int64) (*User, error) {
	accounts, err := p.GetAccounts()
	if err != nil {
		return nil, err
	}

	sid := fmt.Sprintf("%d", id)

	for _, a := range accounts {
		if a.ID == sid {
			return a, nil
		}
	}

	return nil, fae.Errorf("account not found: %s", sid)
}
