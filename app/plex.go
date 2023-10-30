package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

var plex *Plex

func setupPlex() error {
	plex = &Plex{
		PinURL:     "https://plex.tv/api/v2/pins",
		Identifier: "dashotv-web",
		Product:    "DashoTV",
		Device:     "DashoTV (Web)",
	}

	data := url.Values{}
	data.Set("strong", "true")
	data.Set("X-Plex-Client-Identifier", plex.Identifier)
	data.Set("X-Plex-Product", plex.Product)
	plex.data = data

	plex.headers = map[string]string{
		"Accept":                   "application/json",
		"Content-Type":             "application/x-www-form-urlencoded",
		"strong":                   "true",
		"X-Plex-Client-Identifier": plex.Identifier,
		"X-Plex-Product":           plex.Product,
	}

	return nil
}

type Plex struct {
	PinURL     string
	Identifier string
	Product    string
	Device     string
	data       url.Values
	headers    map[string]string
}

func (p *Plex) request(method, url string, values url.Values) (*http.Response, error) {
	r, err := http.NewRequest("POST", p.PinURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	for k, v := range p.headers {
		r.Header.Set(k, v)
	}

	return http.DefaultClient.Do(r)
}

// CreatePin returns a new pin from the plex api
func (p *Plex) CreatePin() (*Pin, error) {
	resp, err := p.request("POST", p.PinURL, p.data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	pin := &Pin{}
	err = json.Unmarshal(d, pin)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response")
	}

	return pin, nil
}

func (p *Plex) CheckPin(pin *Pin) (bool, error) {
	resp, err := p.request("POST", fmt.Sprintf("%s/%d", p.PinURL, pin.Pin), p.data)
	if err != nil {
		return false, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to read response")
	}

	newPin := &Pin{}
	err = json.Unmarshal(d, newPin)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshal response")
	}

	if newPin.Token == "" {
		return false, errors.Errorf("pin not authorized")
	}

	pin.Token = newPin.Token
	pin.Product = newPin.Product
	pin.Identifier = newPin.Identifier

	err = db.Pin.Update(pin)
	if err != nil {
		return false, errors.Wrap(err, "failed to save pin")
	}

	return true, nil
}

func (p *Plex) getAuthUrl(pin *Pin) string {
	base := "https://app.plex.tv/auth/#?"
	data := url.Values{}
	data.Set("clientID", p.Identifier)
	data.Set("code", pin.Code)
	data.Set("forwardUrl", fmt.Sprintf("%s/auth?pin=%d", cfg.Plex, pin.Pin))
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
