package plex

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

func (p *Client) Search(query, section string) ([]SearchMetadata, error) {
	id, err := p.LibraryType(section)
	if err != nil {
		return nil, err
	}

	dest := &PlexSearch{}
	path := fmt.Sprintf("/library/sections/%s/search", section)

	params := url.Values{}
	params.Set("X-Plex-Token", p.Token)
	params.Set("title", query)
	params.Set("type", fmt.Sprintf("%d", id))
	params.Set("limit", "25")
	params.Set("sort", "createdAt:desc")

	resp, err := p._server().SetResult(dest).
		SetQueryParamsFromValues(params).
		Get(path)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get search: %s", resp.Status())
	}

	// app.Log.Debugf("search req url: %s", resp.Request.URL)
	// app.Log.Debugf("search result: %s", resp.String())
	return dest.MediaContainer.Metadata, nil
}

type PlexSearch struct {
	MediaContainer struct {
		Size         int64            `json:"size"`
		SectionID    int64            `json:"sectionID"`
		AllowSync    bool             `json:"allowSync"`
		Art          string           `json:"art"`
		Identifier   string           `json:"identifier"`
		SectionTitle string           `json:"librarySectionTitle"`
		SectionUUID  string           `json:"librarySectionUUID"`
		Title        string           `json:"title1"`
		Subtitle     string           `json:"title2"`
		Metadata     []SearchMetadata `json:"Metadata,omitempty"`
	} `json:"MediaContainer"`
}

type SearchMetadata struct {
	LibrarySectionTitle   string     `json:"librarySectionTitle"`
	Score                 string     `json:"score"`
	RatingKey             string     `json:"ratingKey"`
	Key                   string     `json:"key"`
	GUID                  string     `json:"guid"`
	Slug                  *string    `json:"slug,omitempty"`
	Studio                *string    `json:"studio,omitempty"`
	Type                  string     `json:"type"`
	Title                 string     `json:"title"`
	LibrarySectionID      int64      `json:"librarySectionID"`
	LibrarySectionKey     string     `json:"librarySectionKey"`
	ContentRating         string     `json:"contentRating"`
	Summary               string     `json:"summary"`
	Index                 *int64     `json:"index,omitempty"`
	AudienceRating        *float64   `json:"audienceRating,omitempty"`
	ViewCount             *int64     `json:"viewCount,omitempty"`
	LastViewedAt          *int64     `json:"lastViewedAt,omitempty"`
	Year                  *int64     `json:"year,omitempty"`
	Tagline               *string    `json:"tagline,omitempty"`
	Thumb                 string     `json:"thumb"`
	Art                   string     `json:"art"`
	Theme                 *string    `json:"theme,omitempty"`
	Duration              int64      `json:"duration"`
	OriginallyAvailableAt string     `json:"originallyAvailableAt"`
	LeafCount             *int64     `json:"leafCount,omitempty"`
	ViewedLeafCount       *int64     `json:"viewedLeafCount,omitempty"`
	ChildCount            *int64     `json:"childCount,omitempty"`
	AddedAt               int64      `json:"addedAt"`
	UpdatedAt             int64      `json:"updatedAt"`
	AudienceRatingImage   *string    `json:"audienceRatingImage,omitempty"`
	Genre                 []Country  `json:"Genre,omitempty"`
	Country               []Country  `json:"Country,omitempty"`
	Role                  []Country  `json:"Role,omitempty"`
	Location              []Location `json:"Location,omitempty"`
	SkipCount             *int64     `json:"skipCount,omitempty"`
	SeasonCount           *int64     `json:"seasonCount,omitempty"`
	Field                 []Field    `json:"Field,omitempty"`
	Rating                *float64   `json:"rating,omitempty"`
	PrimaryExtraKey       *string    `json:"primaryExtraKey,omitempty"`
	RatingImage           *string    `json:"ratingImage,omitempty"`
	Media                 []Media    `json:"Media,omitempty"`
	Director              []Country  `json:"Director,omitempty"`
	Writer                []Country  `json:"Writer,omitempty"`
	ChapterSource         *string    `json:"chapterSource,omitempty"`
	ParentRatingKey       *string    `json:"parentRatingKey,omitempty"`
	GrandparentRatingKey  *string    `json:"grandparentRatingKey,omitempty"`
	ParentGUID            *string    `json:"parentGuid,omitempty"`
	GrandparentGUID       *string    `json:"grandparentGuid,omitempty"`
	TitleSort             *string    `json:"titleSort,omitempty"`
	GrandparentKey        *string    `json:"grandparentKey,omitempty"`
	ParentKey             *string    `json:"parentKey,omitempty"`
	GrandparentTitle      *string    `json:"grandparentTitle,omitempty"`
	ParentTitle           *string    `json:"parentTitle,omitempty"`
	OriginalTitle         *string    `json:"originalTitle,omitempty"`
	ParentIndex           *int64     `json:"parentIndex,omitempty"`
	ParentYear            *int64     `json:"parentYear,omitempty"`
	ParentThumb           *string    `json:"parentThumb,omitempty"`
	GrandparentThumb      *string    `json:"grandparentThumb,omitempty"`
	GrandparentArt        *string    `json:"grandparentArt,omitempty"`
	GrandparentTheme      *string    `json:"grandparentTheme,omitempty"`
	GrandparentSlug       *string    `json:"grandparentSlug,omitempty"`
}
