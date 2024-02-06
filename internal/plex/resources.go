package plex

import (
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
)

func (p *Client) GetResources() ([]*Resource, error) {
	out := []*Resource{}

	params := url.Values{}
	// params.Set("includeHttps", "1")
	// params.Set("includeRelay", "1")
	// params.Set("includeLocal", "1")
	// params.Set("includeIPv6", "1")

	resp, err := p._plextv().
		SetHeaders(p.Headers).
		SetQueryParamsFromValues(params).
		Get("/resources")
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get resources: %s", resp.Status())
	}
	if err := json.Unmarshal(resp.Body(), &out); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}
	return out, nil
}

type Resource struct {
	Name             string `json:"name"`
	Product          string `json:"product"`
	ProductVersion   string `json:"productVersion"`
	Platform         string `json:"platform"`
	PlatformVersion  string `json:"platformVersion"`
	Device           string `json:"device"`
	ClientIdentifier string `json:"clientIdentifier"`
	CreatedAt        string `json:"createdAt"`
	LastSeenAt       string `json:"lastSeenAt"`
	Provides         string `json:"provides"`
	// OwnerID                interface{}  `json:"ownerId"`
	// SourceTitle            interface{}  `json:"sourceTitle"`
	PublicAddress          string       `json:"publicAddress"`
	AccessToken            *string      `json:"accessToken"`
	Owned                  bool         `json:"owned"`
	Home                   bool         `json:"home"`
	Synced                 bool         `json:"synced"`
	Relay                  bool         `json:"relay"`
	Presence               bool         `json:"presence"`
	HTTPSRequired          bool         `json:"httpsRequired"`
	PublicAddressMatches   bool         `json:"publicAddressMatches"`
	DNSRebindingProtection *bool        `json:"dnsRebindingProtection,omitempty"`
	NatLoopbackSupported   *bool        `json:"natLoopbackSupported,omitempty"`
	Connections            []Connection `json:"connections"`
}

type Connection struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Port     int64  `json:"port"`
	URI      string `json:"uri"`
	Local    bool   `json:"local"`
	Relay    bool   `json:"relay"`
	IPv6     bool   `json:"IPv6"`
}
