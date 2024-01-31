package plex

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

type Pin struct {
	ID         int    `json:"id"`
	Code       string `json:"code"`
	Product    string `json:"product"`
	Token      string `json:"authToken"`
	Identifier string `json:"clientIdentifier"`
}

// CreatePin returns a new pin from the plex api
func (p *Client) CreatePin() (*Pin, error) {
	pin := &Pin{}
	resp, err := p._plextv().SetResult(pin).SetQueryParamsFromValues(p.data).Post("/pins")
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to create pin: %s", resp.Status())
	}
	return pin, nil
}

func (p *Client) CheckPin(pin *Pin) (bool, error) {
	params := url.Values{}
	params.Set("code", pin.Code)
	params.Set("X-Plex-Client-Identifier", p.ClientIdentifier)

	newPin := &Pin{}
	resp, err := p._plextv().SetResult(newPin).
		SetHeader("code", pin.Code).
		SetQueryParamsFromValues(params).
		Get(fmt.Sprintf("/pins/%d", pin.ID))
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

	// err = app.DB.Pin.Update(pin)
	// if err != nil {
	// 	return false, errors.Wrap(err, "failed to update token")
	// }

	return true, nil
}

func (p *Client) GetAuthUrl(redirect string, pin *Pin) string {
	base := "https://app.plex.tv/auth/#?"
	data := url.Values{}
	data.Set("clientID", p.ClientIdentifier)
	data.Set("code", pin.Code)
	data.Set("forwardUrl", fmt.Sprintf("%s/auth?pin=%d", redirect, pin.ID))
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
