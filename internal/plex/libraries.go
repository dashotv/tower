package plex

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dashotv/fae"
)

/*
	{
	        "allowSync": true,
	        "art": "/:/resources/show-fanart.jpg",
	        "composite": "/library/sections/2/composite/1704818054",
	        "filters": true,
	        "refreshing": true,
	        "thumb": "/:/resources/show.png",
	        "key": "2",
	        "type": "show",
	        "title": "TV Shows",
	        "agent": "tv.plex.agents.series",
	        "scanner": " TV Series",
	        "language": "en-US",
	        "uuid": "e35a54e6-79ff-4cde-98c6-c05c1a09c821",
	        "updatedAt": 1704818112,
	        "createdAt": 1384056048,
	        "scannedAt": 1704818054,
	        "content": true,
	        "directory": true,
	        "contentChangedAt": 41632394,
	        "hidden": 0,
	        "Location": [
	          {
	            "id": 29,
	            "path": "/mnt/media/tv"
	          }
	        ]
	      },
*/
type LibraryLocation struct {
	ID   int64  `json:"id"`
	Path string `json:"path"`
}

type Library struct {
	AllowSync        bool               `json:"allowSync"`
	Art              string             `json:"art"`
	Composite        string             `json:"composite"`
	Filters          bool               `json:"filters"`
	Refreshing       bool               `json:"refreshing"`
	Thumb            string             `json:"thumb"`
	Key              string             `json:"key"`
	Type             string             `json:"type"`
	Title            string             `json:"title"`
	Agent            string             `json:"agent"`
	Scanner          string             `json:"scanner"`
	Language         string             `json:"language"`
	UUID             string             `json:"uuid"`
	UpdatedAt        int64              `json:"updatedAt"`
	CreatedAt        int64              `json:"createdAt"`
	ScannedAt        int64              `json:"scannedAt"`
	Content          bool               `json:"content"`
	Directory        bool               `json:"directory"`
	ContentChangedAt int64              `json:"contentChangedAt"`
	Hidden           int64              `json:"hidden"`
	Locations        []*LibraryLocation `json:"Location"`
}

type Libraries struct {
	MediaContainer struct {
		Size        int64      `json:"size"`
		Directories []*Library `json:"Directory"`
	} `json:"MediaContainer"`
}

func (p *Client) GetLibraries() ([]*Library, error) {
	dest := &Libraries{}
	resp, err := p._server().SetResult(dest).SetFormDataFromValues(p.data).Get("/library/sections")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fae.Errorf("failed to get libraries: %s", resp.Status())
	}

	return dest.MediaContainer.Directories, nil
}

type LibrarySectionResponse struct {
	MediaContainer struct {
		Size        int64      `json:"size"`
		Directories []*Library `json:"Directory"`
	} `json:"MediaContainer"`
}

func (p *Client) GetLibrarySection(section string, directory string, libtype string, start, limit int) ([]*LeavesMetadata, int64, error) {
	dest := &LeavesMetadataContainer{}
	req := p._server().
		SetResult(dest).
		SetHeaders(p.Headers).
		SetHeader("X-Plex-Container-Start", fmt.Sprintf("%d", start)).
		SetHeader("X-Plex-Container-Size", fmt.Sprintf("%d", limit))
	if libtype != "" {
		req.SetQueryParam("type", libtype)
	}
	resp, err := req.Get(fmt.Sprintf("/library/sections/%s/%s", section, directory))
	if err != nil {
		return nil, 0, err
	}
	if !resp.IsSuccess() {
		return nil, 0, fae.Errorf("failed to get libraries: %s", resp.Status())
	}

	return dest.MediaContainer.Metadata, dest.MediaContainer.TotalSize, nil
}

func (p *Client) LibraryType(section string) (int, error) {
	t, err := p.LibraryTypeName(section)
	if err != nil {
		return 0, err
	}
	if t == "" {
		return 0, fae.Errorf("library section %s not found", section)
	}

	id := p.LibraryTypeID(t)
	if id == LibraryTypeUnknown {
		return 0, fae.Errorf("library section %s unknown", section)
	}

	return id, nil
}

func (p *Client) LibraryTypeName(section string) (string, error) {
	resp, err := p.GetLibraries()
	if err != nil {
		return "", err
	}
	for _, r := range resp {
		if r.Key == section {
			return r.Type, nil
		}
	}
	return "", fae.Errorf("library section %s not found", section)
}

func (p *Client) LibraryByPath(path string) (*Library, error) {
	resp, err := p.GetLibraries()
	if err != nil {
		return nil, err
	}
	for _, r := range resp {
		for _, l := range r.Locations {
			if l.Path == path {
				return r, nil
			}
		}
	}
	return nil, fae.Errorf("library path %s not found", path)
}

func (p *Client) RefreshLibraryPath(path string) error {
	l, err := p.LibraryByPath(filepath.Dir(path))
	if err != nil {
		return err
	}

	path = strings.ReplaceAll(path, " ", "+")

	resp, err := p._server().
		SetFormDataFromValues(p.data).
		Get(fmt.Sprintf("/library/sections/%s/refresh?path=%s", l.Key, path))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return fae.Errorf("failed to refresh library: %s", resp.Status())
	}

	return nil
}

func (p *Client) LibraryTypeID(t string) int {
	switch t {
	case "movie":
		return LibraryTypeMovie
	case "show":
		return LibraryTypeShow
	default:
		return LibraryTypeUnknown
	}
}
