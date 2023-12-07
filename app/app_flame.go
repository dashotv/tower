package app

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/dashotv/flame/qbt"
)

var flameClient *Flame

type Flame struct {
	URL string
	c   *resty.Client
}

func setupFlame() error {
	flameClient = &Flame{
		URL: "http://flame:9001",
		c:   resty.New().SetBaseURL("http://flame:9001"),
	}
	return nil
}

type FlameNzbAddResponse struct {
	Error bool `json:"error"`
	ID    int  `json:"id"`
}

func (c *Flame) Torrent(thash string) (*qbt.Torrent, error) {
	res := &qbt.Response{}
	resp, err := c.c.R().
		SetHeader("Accept", "application/json").
		SetResult(res).
		ForceContentType("application/json").
		Get("/qbittorrents/")
	if err != nil {
		return nil, errors.Wrap(err, "failed to load torrent")
	}
	if resp.IsError() {
		return nil, errors.Errorf("failed to load torrent: %s", resp.Status())
	}

	for _, t := range res.Torrents {
		if strings.ToLower(t.Hash) == strings.ToLower(thash) {
			return t, nil
		}
	}

	return nil, errors.Errorf("torrent not found: %s", thash)
}

func (c *Flame) LoadNzb(d *Download, url string) (string, error) {
	enc := base64.StdEncoding.EncodeToString([]byte(url))
	did := d.ID.Hex()
	hash := did[len(did)-4:]
	res := &FlameNzbAddResponse{}

	resp, err := c.c.R().
		SetQueryParam("url", enc).
		SetQueryParam("category", "Series").
		SetQueryParam("name", fmt.Sprintf("[%s] %s %s", hash, d.Medium.Title, d.Medium.Display)).
		SetResult(res).
		Get("/nzbs/add")
	if err != nil {
		return "", errors.Wrap(err, "failed to load nzb")
	}
	if resp.IsError() {
		return "", errors.Errorf("failed to load nzb: %s", resp.Status())
	}
	if res.Error {
		return "", errors.New("failed to load nzb")
	}

	return fmt.Sprintf("%d", res.ID), nil
}

type flameTorrentAddResponse struct {
	Error    bool   `json:"error"`
	Infohash string `json:"infohash"`
}

func (c *Flame) LoadTorrent(d *Download, url string) (string, error) {
	enc := base64.StdEncoding.EncodeToString([]byte(url))
	res := &flameTorrentAddResponse{}
	resp, err := c.c.R().
		SetQueryParam("url", enc).
		SetResult(res).
		Get("/qbittorrents/add")
	if err != nil {
		return "", errors.Wrap(err, "failed to load torrent")
	}
	if resp.IsError() {
		return "", errors.Errorf("failed to load torrent: %s", resp.Status())
	}
	if res.Error {
		return "", errors.New("failed to load torrent")
	}

	return res.Infohash, nil
}