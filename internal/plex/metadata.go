package plex

import (
	"github.com/pkg/errors"
)

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

type PlexLeavesMetadataContainer struct {
	MediaContainer struct {
		Size     int64                 `json:"size"`
		Metadata []*PlexLeavesMetadata `json:"Metadata"`
	} `json:"MediaContainer"`
}
type PlexLeavesMetadata struct {
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
func (p *Client) GetSeriesEpisodes(key string) ([]*PlexLeavesMetadata, error) {
	m := &PlexLeavesMetadataContainer{}
	resp, err := p._server().SetResult(m).SetFormDataFromValues(p.data).Get("/library/metadata/" + key + "/allLeaves")
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get viewed: %s", resp.Status())
	}
	// fmt.Printf("unwatched: %s\n", resp.String())
	return m.MediaContainer.Metadata, nil
}
func (p *Client) GetSeriesEpisodesUnwatched(key string) (*PlexLeavesMetadata, error) {
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
