package app

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/dashotv/flame/metube"
	"github.com/dashotv/flame/qbt"
)

func init() {
	initializers = append(initializers, setupFlame)
}

type Flame struct {
	URL string
	c   *resty.Client
}

func setupFlame(app *Application) error {
	app.Flame = &Flame{
		c: resty.New().SetBaseURL(app.Config.FlameURL),
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

func (c *Flame) LoadTorrent(_ *Download, url string) (string, error) {
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

func (c *Flame) RemoveTorrent(thash string) error {
	resp, err := c.c.R().
		SetQueryParam("infohash", thash).
		Get("/qbittorrents/remove")
	if err != nil {
		return errors.Wrap(err, "failed to remove torrent")
	}
	if resp.IsError() {
		return errors.Errorf("failed to remove torrent: %s", resp.Status())
	}

	app.Log.Debugf("Flame::RemoveTorrent: %s", resp.Body())
	return nil
}

type MetubeAddResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (c *Flame) LoadMetube(name string, url string, autoStart bool) error {
	app.Log.Named("flame").Debugf("LoadMetube: %s %s %t", name, url, autoStart)
	enc := base64.StdEncoding.EncodeToString([]byte(url))
	res := &MetubeAddResponse{}
	resp, err := c.c.R().
		SetQueryParam("url", enc).
		SetQueryParam("name", name).
		SetQueryParam("auto_start", fmt.Sprintf("%t", autoStart)).
		SetResult(res).
		Get("/metube/add")
	if err != nil {
		return errors.Wrap(err, "failed to load metube")
	}
	if resp.IsError() {
		return errors.Errorf("failed to load metube: %s", resp.Status())
	}
	if res.Error {
		return errors.Errorf("failed to load metube: %s", res.Message)
	}

	return nil
}

type MetubeHistory struct {
	Error   bool                    `json:"error"`
	History *metube.HistoryResponse `json:"history"`
}

func (c *Flame) MetubeHistory() (*metube.HistoryResponse, error) {
	res := &MetubeHistory{}
	resp, err := c.c.R().
		SetResult(res).
		SetHeader("Accept", "application/json").
		Get("/metube/")
	if err != nil {
		return nil, errors.Wrap(err, "failed to load torrent")
	}
	if resp.IsError() {
		return nil, errors.Errorf("failed to load torrent: %s", resp.Status())
	}
	if res.Error {
		return nil, errors.New("failed to load metube")
	}

	return res.History, nil
}
