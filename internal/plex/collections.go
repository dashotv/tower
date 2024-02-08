package plex

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type PlexCollectionCreate struct {
	MediaContainer struct {
		Directory []struct {
			RatingKey string `json:"ratingKey"`
			Key       string `json:"key"`
			Guid      string `json:"guid"`
			Type      string `json:"type"`
			Title     string `json:"title"`
			Subtype   string `json:"subtype"`
			Summary   string `json:"summary"`
			Thumb     string `json:"thumb"`
			AddedAt   int64  `json:"addedAt"`
			UpdatedAt int64  `json:"updatedAt"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}

func (p *Client) CreateCollection(title, section, firstKey string) (*PlexCollectionCreate, error) {
	id, err := p.LibraryType(section)
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("X-Plex-Token", p.Token)
	data.Set("title", title)
	data.Set("sectionId", section)
	data.Set("type", fmt.Sprintf("%d", id))
	data.Set("smart", "0")
	data.Set("uri", fmt.Sprintf("server://%s/com.plexapp.plugins.library/library/metadata/%s", p.MachineIdentifier, firstKey))

	dest := &PlexCollectionCreate{}
	resp, err := p._server().
		SetResult(dest).
		SetQueryParamsFromValues(data).
		Post("/library/collections")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to create collection: %s", resp.Status())
	}

	return dest, nil
}

type PlexLibrariesCollectionResponse struct {
	MediaContainer struct {
		Size         int64             `json:"size"`
		AllowSync    bool              `json:"allowSync"`
		Identifier   string            `json:"identifier"`
		LibraryID    int64             `json:"librarySectionID"`
		LibraryTitle string            `json:"librarySectionTitle"`
		LibraryUUID  string            `json:"librarySectionUUID"`
		Title        string            `json:"title1"`
		Subtitle     string            `json:"title2"`
		Metadata     []*PlexCollection `json:"Metadata,omitempty"`
	} `json:"MediaContainer"`
}

type PlexCollectionResponse struct {
	MediaContainer struct {
		Size         int64             `json:"size"`
		AllowSync    bool              `json:"allowSync"`
		Identifier   string            `json:"identifier"`
		LibraryID    int64             `json:"librarySectionID"`
		LibraryTitle string            `json:"librarySectionTitle"`
		LibraryUUID  string            `json:"librarySectionUUID"`
		Title        string            `json:"title1"`
		Subtitle     string            `json:"title2"`
		Directory    []*PlexCollection `json:"Metadata,omitempty"`
	} `json:"MediaContainer"`
}
type PlexCollection struct {
	RatingKey    string                 `json:"ratingKey"`
	Key          string                 `json:"key"`
	GUID         string                 `json:"guid"`
	Type         string                 `json:"type"`
	Title        string                 `json:"title"`
	LibraryID    int64                  `json:"librarySectionID"`
	LibraryTitle string                 `json:"librarySectionTitle"`
	LibraryKey   string                 `json:"librarySectionKey"`
	Subtype      string                 `json:"subtype"`
	Summary      string                 `json:"summary"`
	Thumb        string                 `json:"thumb"`
	AddedAt      int64                  `json:"addedAt"`
	UpdatedAt    int64                  `json:"updatedAt"`
	ChildCount   string                 `json:"childCount"`
	MaxYear      string                 `json:"maxYear"`
	MinYear      string                 `json:"minYear"`
	Children     []*PlexCollectionChild `json:"children,omitempty"`
}
type PlexCollectionChildrenResponse struct {
	MediaContainer struct {
		Size      int64                  `json:"size"`
		Directory []*PlexCollectionChild `json:"Metadata,omitempty"`
	} `json:"MediaContainer"`
}
type PlexCollectionChild struct {
	RatingKey    string `json:"ratingKey"`
	Key          string `json:"key"`
	GUID         string `json:"guid"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	LibraryID    int64  `json:"librarySectionID"`
	LibraryTitle string `json:"librarySectionTitle"`
	LibraryKey   string `json:"librarySectionKey"`
	Summary      string `json:"summary"`
	Thumb        string `json:"thumb"`
	AddedAt      int64  `json:"addedAt"`
	UpdatedAt    int64  `json:"updatedAt"`
}

func (p *Client) DeleteCollection(ratingKey string) error {
	data := url.Values{}
	data.Set("X-Plex-Token", p.Token)

	resp, err := p._server().
		SetQueryParamsFromValues(data).
		Delete(fmt.Sprintf("/library/collections/%s", ratingKey))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return errors.Errorf("failed to delete collection: %s", resp.Status())
	}

	return nil
}

func (p *Client) ListCollections(section string) ([]*PlexCollection, error) {
	data := url.Values{}
	data.Set("X-Plex-Token", p.Token)

	dest := &PlexLibrariesCollectionResponse{}
	resp, err := p._server().
		SetResult(dest).
		SetQueryParamsFromValues(p.data).
		Get("/library/sections/" + section + "/collections")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get collections: %s", resp.Status())
	}

	return dest.MediaContainer.Metadata, nil
}

func (p *Client) GetCollection(ratingKey string) (*PlexCollection, error) {
	data := url.Values{}
	data.Set("X-Plex-Token", p.Token)

	dest := &PlexCollectionResponse{}
	resp, err := p._server().
		SetResult(dest).
		SetQueryParamsFromValues(p.data).
		Get("/library/collections/" + ratingKey)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to update collection: %s", resp.Status())
	}
	if len(dest.MediaContainer.Directory) != 1 {
		return nil, errors.Errorf("api response found %d directories, wanted 1", len(dest.MediaContainer.Directory))
	}
	children, err := p.GetCollectionChildren(ratingKey)
	if err != nil {
		return nil, err
	}

	r := dest.MediaContainer.Directory[0]
	r.Children = children

	return r, nil
}

func (p *Client) GetCollectionChildren(ratingKey string) ([]*PlexCollectionChild, error) {
	data := url.Values{}
	data.Set("X-Plex-Token", p.Token)

	dest := &PlexCollectionChildrenResponse{}
	resp, err := p._server().
		SetResult(dest).
		SetQueryParamsFromValues(p.data).
		Get("/library/collections/" + ratingKey + "/children")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get collection children: %s", resp.Status())
	}
	// fmt.Printf("collection: %s\n", resp.String())
	return dest.MediaContainer.Directory, nil
}

func (p *Client) UpdateCollection(section, ratingKey string, keys []string) error {
	existing, err := p.GetCollection(ratingKey)
	if err != nil {
		return err
	}

	existingKeys := lo.Map(existing.Children, func(c *PlexCollectionChild, i int) string {
		return c.RatingKey
	})

	add, remove := lo.Difference(keys, existingKeys)
	if len(add) > 0 {
		for _, k := range add {
			if err := p.addCollectionItem(ratingKey, k); err != nil {
				return err
			}
		}
	}
	if len(remove) > 0 {
		for _, k := range remove {
			if err := p.removeCollectionItem(ratingKey, k); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Client) addCollectionItem(ratingKey, newKey string) error {
	data := url.Values{}
	data.Set("uri", fmt.Sprintf("server://%s/com.plexapp.plugins.library/library/metadata/%s", p.MachineIdentifier, newKey))

	resp, err := p._server().
		SetHeaders(p.Headers).
		SetQueryParamsFromValues(data).
		Put("/library/collections/" + ratingKey + "/items")
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return errors.Errorf("failed to add to collection: %s", resp.Status())
	}

	return nil
}

func (p *Client) removeCollectionItem(ratingKey, rmKey string) error {
	data := url.Values{}
	data.Set("X-Plex-Token", p.Token)
	data.Set("excludeAllLeaves", "1")

	resp, err := p._server().
		SetQueryParamsFromValues(data).
		Delete(fmt.Sprintf("/library/collections/%s/children/%s", ratingKey, rmKey))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return errors.Errorf("failed to remove from collection: %s", resp.Status())
	}

	return nil
}

// func (p *Plex) GetCollections(token string) (*PlexCollections, error) {
// 	dest := &PlexCollections{}
// 	resp, err := p.metadata().SetResult(dest).SetHeader("X-Plex-Token", token).Get("/library/sections")
// 	if err != nil {
// 		return dest, err
// 	}
// 	if !resp.IsSuccess() {
// 		return dest, errors.Errorf("failed to get collections: %s", resp.Status())
// 	}
//
// 	return dest, nil
// }
