package plex

import (
	"fmt"

	"github.com/pkg/errors"
)

type WatchlistOpts struct {
	Filter string // all, or ???
	Sort   string // ???
	Type   string // library type? movie, show, episode, artist, album, track?
}

func (p *Client) GetWatchlist(token string) (*PlexWatchlist, error) {
	dest := &PlexWatchlist{}
	opts := &WatchlistOpts{Filter: "all"}
	u := fmt.Sprintf("/library/sections/watchlist/%s", opts.Filter)
	resp, err := p._metadata().SetResult(dest).SetHeader("X-Plex-Token", token).Get(u)
	if err != nil {
		return dest, err
	}
	if !resp.IsSuccess() {
		return dest, errors.Errorf("failed to get watchlist: %s", resp.Status())
	}

	return dest, nil
}

func (p *Client) GetWatchlistDetail(token string, w *PlexWatchlist) ([]*PlexWatchlistDetail, error) {
	out := []*PlexWatchlistDetail{}
	for _, d := range w.MediaContainer.Metadata {
		dest := &PlexWatchlistDetail{}
		resp, err := p._metadata().
			SetResult(dest).
			SetHeader("X-Plex-Token", token).
			Get("/library/metadata/" + d.RatingKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to make watchlistdetails request")
		}
		if !resp.IsSuccess() {
			return nil, errors.Errorf("failed to get watchlistdetails: %s", resp.Status())
		}
		out = append(out, dest)
	}
	return out, nil
}

type PlexWatchlist struct {
	MediaContainer struct {
		LibrarySectionID    string `json:"librarySectionID"`
		LibrarySectionTitle string `json:"librarySectionTitle"`
		Offset              int64  `json:"offset"`
		TotalSize           int64  `json:"totalSize"`
		Identifier          string `json:"identifier"`
		Size                int64  `json:"size"`
		Metadata            []struct {
			Art                   string   `json:"art"`
			Banner                string   `json:"banner"`
			GUID                  string   `json:"guid"`
			Key                   string   `json:"key"`
			Rating                float64  `json:"rating"`
			RatingKey             string   `json:"ratingKey"`
			Studio                string   `json:"studio"`
			Type                  string   `json:"type"`
			Theme                 *string  `json:"theme,omitempty"`
			Thumb                 string   `json:"thumb"`
			AddedAt               int64    `json:"addedAt"`
			Duration              int64    `json:"duration"`
			PublicPagesURL        string   `json:"publicPagesURL"`
			Slug                  string   `json:"slug"`
			UserState             bool     `json:"userState"`
			Title                 string   `json:"title"`
			OriginalTitle         *string  `json:"originalTitle,omitempty"`
			LeafCount             *int64   `json:"leafCount,omitempty"`
			ChildCount            *int64   `json:"childCount,omitempty"`
			ContentRating         string   `json:"contentRating"`
			OriginallyAvailableAt string   `json:"originallyAvailableAt"`
			Year                  int64    `json:"year"`
			RatingImage           string   `json:"ratingImage"`
			ImdbRatingCount       int64    `json:"imdbRatingCount"`
			Image                 []Image  `json:"Image"`
			PrimaryExtraKey       *string  `json:"primaryExtraKey,omitempty"`
			IsContinuingSeries    *bool    `json:"isContinuingSeries,omitempty"`
			Tagline               *string  `json:"tagline,omitempty"`
			AudienceRating        *float64 `json:"audienceRating,omitempty"`
			AudienceRatingImage   *string  `json:"audienceRatingImage,omitempty"`
			Subtype               *string  `json:"subtype,omitempty"`
			SkipChildren          *bool    `json:"skipChildren,omitempty"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}

type PlexWatchlistDetail struct {
	MediaContainer struct {
		Offset     int64  `json:"offset"`
		TotalSize  int64  `json:"totalSize"`
		Identifier string `json:"identifier"`
		Size       int64  `json:"size"`
		Metadata   []struct {
			Art                   string     `json:"art"`
			Banner                string     `json:"banner"`
			MetadatumGUID         string     `json:"guid"`
			Key                   string     `json:"key"`
			MetadatumRating       float64    `json:"rating"`
			RatingKey             string     `json:"ratingKey"`
			MetadatumStudio       string     `json:"studio"`
			Summary               string     `json:"summary"`
			Tagline               string     `json:"tagline"`
			Type                  string     `json:"type"`
			Theme                 string     `json:"theme"`
			Thumb                 string     `json:"thumb"`
			AddedAt               int64      `json:"addedAt"`
			Duration              int64      `json:"duration"`
			PublicPagesURL        string     `json:"publicPagesURL"`
			Slug                  string     `json:"slug"`
			UserState             bool       `json:"userState"`
			Title                 string     `json:"title"`
			LeafCount             int64      `json:"leafCount"`
			ChildCount            int64      `json:"childCount"`
			ContentRating         string     `json:"contentRating"`
			OriginallyAvailableAt string     `json:"originallyAvailableAt"`
			Year                  int64      `json:"year"`
			RatingImage           string     `json:"ratingImage"`
			ImdbRatingCount       int64      `json:"imdbRatingCount"`
			Image                 []Image    `json:"Image"`
			Genre                 []Genre    `json:"Genre"`
			GUID                  []GUID     `json:"Guid"`
			Country               []Country  `json:"Country"`
			Role                  []Role     `json:"Role"`
			Director              []Director `json:"Director"`
			Producer              []Director `json:"Producer"`
			Writer                []Director `json:"Writer"`
			Network               []Country  `json:"Network"`
			Rating                []Rating   `json:"Rating"`
			Similar               []Similar  `json:"Similar"`
			Studio                []Country  `json:"Studio"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}
