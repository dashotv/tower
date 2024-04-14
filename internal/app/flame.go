package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/dashotv/fae"
	flame "github.com/dashotv/flame/client"
	"github.com/dashotv/flame/metube"
	"github.com/dashotv/flame/qbt"
)

func init() {
	initializers = append(initializers, setupFlame)
}

func setupFlame(app *Application) error {
	app.Flame = flame.New(app.Config.FlameURL)
	return nil
}

func (a *Application) FlameTorrent(thash string) (*qbt.Torrent, error) {
	resp, err := a.Flame.Qbittorrents.Index(context.Background())
	if err != nil {
		return nil, fae.Wrap(err, "failed to load torrents")
	}
	for _, t := range resp.Result.Torrents {
		if strings.ToLower(t.Hash) == strings.ToLower(thash) {
			return t, nil
		}
	}
	return nil, fae.Errorf("torrent not found: %s", thash)
}

func (a *Application) FlameMetubeHistory() (*metube.HistoryResponse, error) {
	resp, err := a.Flame.Metube.Index(context.Background())
	if err != nil {
		return nil, fae.Wrap(err, "failed to load metube history")
	}
	return resp.Result, nil
}

func (a *Application) FlameAdd(d *Download) (string, error) {
	a.Log.Named("flame").Debugf("FlameAdd: %s - %s :: %s", d.Title, d.Display, d.URL)
	if d.IsNzb() {
		return a.FlameNzbAdd(d)
	}
	if d.IsTorrent() {
		return a.FlameTorrentAdd(d)
	}
	if d.IsMetube() {
		return a.FlameMetubeAdd(d)
	}
	return "", fae.Errorf("unsupported download type: %s", d.URL)
}

func (a *Application) FlameNzbAdd(d *Download) (string, error) {
	if d.Medium == nil {
		return "", fae.New("missing medium")
	}

	url, err := d.GetURL()
	if err != nil {
		return "", fae.Wrap(err, "getting url")
	}
	enc := base64.StdEncoding.EncodeToString([]byte(url))
	did := d.ID.Hex()
	hash := did[len(did)-4:]
	name := fmt.Sprintf("[%s] %s %s", hash, d.Title, d.Display)
	category := "Series"
	if d.Medium.Type == "Movie" {
		category = "Movies"
	}

	req := &flame.NzbsAddRequest{
		URL:      enc,
		Name:     name,
		Category: category,
	}

	resp, err := a.Flame.Nzbs.Add(context.Background(), req)
	if err != nil {
		return "", fae.Wrap(err, "failed to load nzb")
	}
	if resp.Error {
		return "", fae.New("failed to load nzb")
	}
	return fmt.Sprintf("%d", resp.Result), nil
}

func (a *Application) FlameMetubeAdd(d *Download) (string, error) {
	url, err := d.GetURL()
	if err != nil {
		return "", fae.Wrap(err, "getting url")
	}
	url = strings.Replace(url, "metube://", "", 1)
	enc := base64.StdEncoding.EncodeToString([]byte(url))
	did := d.ID.Hex()

	req := &flame.MetubeAddRequest{
		URL:       enc,
		Name:      did,
		AutoStart: true,
	}

	resp, err := a.Flame.Metube.Add(context.Background(), req)
	if err != nil {
		return "", fae.Wrap(err, "failed to load metube")
	}
	if resp.Error {
		return "", fae.New(resp.Message)
	}

	return "M", nil
}
func (a *Application) FlameTorrentAdd(d *Download) (string, error) {
	url, err := d.GetURL()
	if err != nil {
		return "", fae.Wrap(err, "getting url")
	}
	enc := base64.StdEncoding.EncodeToString([]byte(url))

	req := &flame.QbittorrentsAddRequest{
		URL: enc,
	}

	resp, err := a.Flame.Qbittorrents.Add(context.Background(), req)
	if err != nil {
		return "", fae.Wrap(err, "failed to load torrent")
	}
	if resp.Error {
		return "", fae.New(resp.Message)
	}

	return strings.ToLower(resp.Result), nil
}
func (a *Application) FlameTorrentRemove(thash string) error {
	req := &flame.QbittorrentsRemoveRequest{
		Infohash: thash,
	}

	resp, err := a.Flame.Qbittorrents.Remove(context.Background(), req)
	if err != nil {
		return fae.Wrap(err, "failed to remove torrent")
	}
	if resp.Error {
		return fae.New(resp.Message)
	}

	return nil
}

type Flame struct {
	URL string
	c   *resty.Client
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
		return nil, fae.Wrap(err, "failed to load torrent")
	}
	if resp.IsError() {
		return nil, fae.Errorf("failed to load torrent: %s", resp.Status())
	}

	for _, t := range res.Torrents {
		if strings.ToLower(t.Hash) == strings.ToLower(thash) {
			return t, nil
		}
	}

	return nil, fae.Errorf("torrent not found: %s", thash)
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
		return "", fae.Wrap(err, "failed to load nzb")
	}
	if resp.IsError() {
		return "", fae.Errorf("failed to load nzb: %s", resp.Status())
	}
	if res.Error {
		return "", fae.New("failed to load nzb")
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
		return "", fae.Wrap(err, "failed to load torrent")
	}
	if resp.IsError() {
		return "", fae.Errorf("failed to load torrent: %s", resp.Status())
	}
	if res.Error {
		return "", fae.New("failed to load torrent")
	}

	return res.Infohash, nil
}

func (c *Flame) RemoveTorrent(thash string) error {
	resp, err := c.c.R().
		SetQueryParam("infohash", thash).
		Get("/qbittorrents/remove")
	if err != nil {
		return fae.Wrap(err, "failed to remove torrent")
	}
	if resp.IsError() {
		return fae.Errorf("failed to remove torrent: %s", resp.Status())
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
		return fae.Wrap(err, "failed to load metube")
	}
	if resp.IsError() {
		return fae.Errorf("failed to load metube: %s", resp.Status())
	}
	if res.Error {
		return fae.Errorf("failed to load metube: %s", res.Message)
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
		return nil, fae.Wrap(err, "failed to load history")
	}
	if resp.IsError() {
		return nil, fae.Errorf("failed to load history: %s", resp.Status())
	}
	if res.Error {
		return nil, fae.New("failed to load history")
	}

	return res.History, nil
}
