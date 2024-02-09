package plex

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

func (p *Client) Play(ratingKey, player string) error {
	queue, err := p.playCreateQueue(ratingKey)
	if err != nil {
		return errors.Wrap(err, "failed to create queue")
	}
	if queue == nil {
		return errors.New("failed to create queue")
	}

	return p.playQueue(queue.MediaContainer.ID, ratingKey, player)
}

func (p *Client) Stop(session string) error {
	params := url.Values{}
	params.Set("sessionId", session)
	params.Set("reason", "")

	resp, err := p._server().
		SetHeaders(p.Headers).
		SetQueryParamsFromValues(params).
		Get("/status/sessions/terminate")
	if err != nil {
		return errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return errors.Errorf("failed to play: %s", resp.Status())
	}

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

func (p *Client) playQueue(queueID int64, ratingKey, player string) error {
	params := url.Values{}
	params.Set("protocol", "http")
	params.Set("address", "10.0.4.61")
	params.Set("port", "32400")
	params.Set("offset", "0")
	params.Set("commandID", "1")
	params.Set("machineIdentifier", p.MachineIdentifier)
	params.Set("type", "video")
	params.Set("containerKey", fmt.Sprintf("/playQueues/%d", queueID))
	params.Set("key", "/library/metadata/"+ratingKey)
	resp, err := p._server().
		// SetDebug(true).
		SetHeaders(p.Headers).
		SetHeader("X-Plex-Target-Client-Identifier", player).
		SetQueryParamsFromValues(params).
		Get("/player/playback/playMedia")
	if err != nil {
		return errors.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return errors.Errorf("failed to play: %s: %s", resp.Status(), resp.String())
	}
	return nil
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
	// fmt.Printf("queue: %s\n", resp.String())
	// fmt.Printf("queue: %+v\n", q)
	return q, nil
}
