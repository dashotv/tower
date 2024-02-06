package plex

import (
	"fmt"
	"net/url"

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
func (p *Client) Play(player, ratingKey string) error {
	queue, err := p.playCreateQueue(ratingKey)
	if err != nil {
		return errors.Wrap(err, "failed to create queue")
	}
	if queue == nil {
		return errors.New("failed to create queue")
	}

	params := url.Values{}
	params.Set("protocol", "http")
	params.Set("address", "10.0.4.61")
	params.Set("port", "32400")
	params.Set("offset", "0")
	params.Set("commandID", "1")
	params.Set("machineIdentifier", p.MachineIdentifier)
	params.Set("type", "video")
	params.Set("containerKey", fmt.Sprintf("/playQueues/%d", queue.MediaContainer.ID))
	params.Set("key", "/library/metadata/"+ratingKey)
	resp, err := p._server().SetHeaders(p.Headers).
		SetHeader("X-Plex-Target-Client-Identifier", player).
		SetQueryParamsFromValues(params).
		Get("/player/playback/playMedia")
	if err != nil {
		return errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return errors.Errorf("failed to play: %s", resp.Status())
	}
	fmt.Printf("play: %s\n", resp.String())
	return nil
}

type PlexQueue struct {
	MediaContainer struct {
		Size            int64  `json:"size"`
		Identifier      string `json:"identifier"`
		MediaTagPrefix  string `json:"mediaTagPrefix"`
		MediaTagVersion int64  `json:"mediaTagVersion"`
		Shuffled        bool   `json:"playQueueShuffled"`
		Source          string `json:"playQueueSourceURI"`
		Version         int64  `json:"playQueueVersion"`
		ID              int64  `json:"playQueueID"`
	}
}

func (p *Client) playCreateQueue(ratingKey string) (*PlexQueue, error) {
	q := &PlexQueue{}

	params := url.Values{}
	params.Set("type", "video")
	params.Set("shuffle", "0")
	params.Set("repeat", "0")
	params.Set("continuous", "0")
	params.Set("own", "1")
	params.Set("uri", fmt.Sprintf("server://%s/com.plexapp.plugins.library/library/metadata/%s", p.MachineIdentifier, ratingKey))

	resp, err := p._server().SetResult(q).SetHeaders(p.Headers).SetQueryParamsFromValues(params).Post("/playQueues")
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to play: %s", resp.Status())
	}
	fmt.Printf("queue: %s\n", resp.String())
	fmt.Printf("queue: %+v\n", q)
	return q, nil
}
