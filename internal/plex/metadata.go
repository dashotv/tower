package plex

import (
	"github.com/dashotv/fae"
)

type LibraryMetadataContainer struct {
	MediaContainer struct {
		Size     int64              `json:"size"`
		Metadata []*LibraryMetadata `json:"Metadata"`
	} `json:"MediaContainer"`
}
type LibraryMetadata struct {
	Key          string `json:"key"`
	RatingKey    string `json:"ratingKey"`
	Leaves       int    `json:"leafCount"`
	Viewed       int    `json:"viewedLeafCount"`
	LastViewedAt int64  `json:"lastViewedAt"`
}

type LeavesMetadataContainer struct {
	MediaContainer struct {
		Size     int64             `json:"size"`
		Metadata []*LeavesMetadata `json:"Metadata"`
	} `json:"MediaContainer"`
}
type LeavesMetadata struct {
	Key              string   `json:"key"`
	RatingKey        string   `json:"ratingKey"`
	LastViewedAt     int64    `json:"lastViewedAt"`
	Title            string   `json:"title"`
	ParentTitle      string   `json:"parentTitle"`
	GrandparentTitle string   `json:"grandparentTitle"`
	Index            int      `json:"index"`
	AddedAt          int64    `json:"addedAt"`
	UpdatedAt        int64    `json:"updatedAt"`
	Media            []*Media `json:"Media"`
}

func (p *Client) GetMetadataByKey(key string) (string, error) {
	resp, err := p._server().SetFormDataFromValues(p.data).Get("/library/metadata/" + key)
	if err != nil {
		return "", fae.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return "", fae.Errorf("failed to get metadata: %s", resp.Status())
	}
	return resp.String(), nil
}
func (p *Client) GetViewedByKey(key string) (*LibraryMetadata, error) {
	m := &LibraryMetadataContainer{}
	resp, err := p._server().SetResult(m).SetFormDataFromValues(p.data).Get("/library/metadata/" + key)
	if err != nil {
		return nil, fae.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, fae.Errorf("failed to get viewed: %s", resp.Status())
	}
	return m.MediaContainer.Metadata[0], nil
}
func (p *Client) GetSeriesEpisodes(key string) ([]*LeavesMetadata, error) {
	m := &LeavesMetadataContainer{}
	resp, err := p._server().SetResult(m).SetFormDataFromValues(p.data).Get("/library/metadata/" + key + "/allLeaves")
	if err != nil {
		return nil, fae.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, fae.Errorf("failed to get viewed: %s", resp.Status())
	}
	// fmt.Printf("unwatched: %s\n", resp.String())
	return m.MediaContainer.Metadata, nil
}
func (p *Client) GetSeriesEpisodesUnwatched(key string) (*LeavesMetadata, error) {
	list, err := p.GetSeriesEpisodes(key)
	if err != nil {
		return nil, err
	}
	for _, ep := range list {
		if ep.LastViewedAt == 0 {
			return ep, nil
		}
	}
	return nil, nil
}
