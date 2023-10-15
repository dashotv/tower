package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var pinUrl = "https://plex.tv/api/v2/pins"
var identifier = "dashotv-web"
var product = "DashoTV"
var device = "DashoTV (Web)"

func PlexIndex(c *gin.Context) {
	// get pin
	pin, err := plexCreatePin()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	err = db.Pin.Save(pin)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	authUrl := plexAuthUrl(pin)
	// c.JSON(200, gin.H{"pin": pin, "authUrl": authUrl})
	c.Redirect(302, authUrl)
}

func PlexAuth(c *gin.Context) {
	id := c.Query("pin")
	pinId, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	list, err := db.Pin.Query().Where("pin", int64(pinId)).Run()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	server.Log.Infof("list: %d %+v", pinId, list)
	if len(list) != 1 {
		c.AbortWithStatusJSON(404, gin.H{"error": "pin not found"})
		return
	}

	pin := list[0]
	_, err = plexCheckPin(pin)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	// TODO: get user from token (call myplex api), maybe background this?
	c.String(200, "Authorization complete!")
}

// plexCreatePin returns a new pin from the plex api
func plexCreatePin() (*Pin, error) {
	data := url.Values{}
	data.Set("strong", "true")
	data.Set("X-Plex-Client-Identifier", identifier)
	data.Set("X-Plex-Product", product)

	r, err := http.NewRequest("POST", pinUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("strong", "true")
	r.Header.Add("X-Plex-Client-Identifier", identifier)
	r.Header.Add("X-Plex-Product", product)

	resp, err := http.DefaultClient.Do(r)
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

func plexCheckPin(pin *Pin) (bool, error) {
	data := url.Values{}
	data.Set("code", pin.Code)
	data.Set("strong", "true")
	data.Set("X-Plex-Client-Identifier", identifier)
	data.Set("X-Plex-Product", product)

	r, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", pinUrl, pin.Pin), strings.NewReader(data.Encode()))
	if err != nil {
		return false, errors.Wrap(err, "failed to create request")
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("strong", "true")
	r.Header.Add("X-Plex-Client-Identifier", identifier)
	r.Header.Add("X-Plex-Product", product)

	resp, err := http.DefaultClient.Do(r)
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

func plexAuthUrl(pin *Pin) string {
	base := "https://app.plex.tv/auth/#?"
	data := url.Values{}
	data.Set("clientID", identifier)
	data.Set("code", pin.Code)
	data.Set("forwardUrl", fmt.Sprintf("%s/auth?pin=%d", cfg.Plex, pin.Pin))
	data.Set("context[device][product]", product)
	data.Set("context[device][version]", "0.1.0")
	data.Set("context[device][deviceName]", device)
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
