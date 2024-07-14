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
		Size      int64             `json:"size"`
		TotalSize int64             `json:"totalSize"`
		Metadata  []*LeavesMetadata `json:"Metadata"`
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
type KeyResponse struct {
	MediaContainer struct {
		Size                int64                  `json:"size"`
		AllowSync           bool                   `json:"allowSync"`
		Identifier          string                 `json:"identifier"`
		LibrarySectionID    int64                  `json:"librarySectionID"`
		LibrarySectionTitle string                 `json:"librarySectionTitle"`
		LibrarySectionUUID  string                 `json:"librarySectionUUID"`
		MediaTagPrefix      string                 `json:"mediaTagPrefix"`
		MediaTagVersion     int64                  `json:"mediaTagVersion"`
		Metadata            []*KeyResponseMetadata `json:"Metadata"`
	} `json:"MediaContainer"`
}

type KeyResponseMetadata struct {
	RatingKey             string   `json:"ratingKey"`
	Key                   string   `json:"key"`
	ParentRatingKey       string   `json:"parentRatingKey"`
	GrandparentRatingKey  string   `json:"grandparentRatingKey"`
	MetadatumGUID         string   `json:"guid"`
	ParentGUID            string   `json:"parentGuid"`
	GrandparentGUID       string   `json:"grandparentGuid"`
	GrandparentSlug       string   `json:"grandparentSlug"`
	Type                  string   `json:"type"`
	Title                 string   `json:"title"`
	GrandparentKey        string   `json:"grandparentKey"`
	ParentKey             string   `json:"parentKey"`
	LibrarySectionTitle   string   `json:"librarySectionTitle"`
	LibrarySectionID      int64    `json:"librarySectionID"`
	LibrarySectionKey     string   `json:"librarySectionKey"`
	GrandparentTitle      string   `json:"grandparentTitle"`
	ParentTitle           string   `json:"parentTitle"`
	Summary               string   `json:"summary"`
	Index                 int64    `json:"index"`
	ParentIndex           int64    `json:"parentIndex"`
	SkipCount             int64    `json:"skipCount"`
	Year                  int64    `json:"year"`
	Thumb                 string   `json:"thumb"`
	Art                   string   `json:"art"`
	ParentThumb           string   `json:"parentThumb"`
	GrandparentThumb      string   `json:"grandparentThumb"`
	GrandparentArt        string   `json:"grandparentArt"`
	Duration              int64    `json:"duration"`
	OriginallyAvailableAt string   `json:"originallyAvailableAt"`
	AddedAt               int64    `json:"addedAt"`
	UpdatedAt             int64    `json:"updatedAt"`
	Media                 []*Media `json:"Media"`
	GUID                  []*GUID  `json:"Guid"`
}

func (p *Client) GetMetadataByKey(key string) ([]*KeyResponseMetadata, error) {
	m := &KeyResponse{}
	params := map[string]string{
		"includeExternalMedia": "1",
		"includePreferences":   "1",
		"skipRefresh":          "1",
	}
	resp, err := p._server().SetResult(m).SetFormDataFromValues(p.data).SetQueryParams(params).Get("/library/metadata/" + key)
	if err != nil {
		return nil, fae.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return nil, fae.Errorf("failed to get metadata: %s", resp.Status())
	}
	return m.MediaContainer.Metadata, nil
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
	resp, err := p._server().SetResult(m).SetHeader("X-Plex-Container-Size", "500").SetFormDataFromValues(p.data).Get("/library/metadata/" + key + "/allLeaves")
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
func (p *Client) PutMetadataPrefs(key string, prefs map[string]string) error {
	resp, err := p._server().SetHeaders(p.Headers).SetQueryParams(prefs).Put("/library/metadata/" + key + "/prefs")
	if err != nil {
		return fae.Wrap(err, "failed to make request")
	}
	if !resp.IsSuccess() {
		return fae.Errorf("failed to put metadata prefs: %s", resp.Status())
	}
	return nil
}
