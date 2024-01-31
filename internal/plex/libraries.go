package plex

import "github.com/pkg/errors"

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
	        "scanner": "Plex TV Series",
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
type PlexLibraryLocation struct {
	ID   int64  `json:"id"`
	Path string `json:"path"`
}

type PlexLibrary struct {
	AllowSync        bool                   `json:"allowSync"`
	Art              string                 `json:"art"`
	Composite        string                 `json:"composite"`
	Filters          bool                   `json:"filters"`
	Refreshing       bool                   `json:"refreshing"`
	Thumb            string                 `json:"thumb"`
	Key              string                 `json:"key"`
	Type             string                 `json:"type"`
	Title            string                 `json:"title"`
	Agent            string                 `json:"agent"`
	Scanner          string                 `json:"scanner"`
	Language         string                 `json:"language"`
	UUID             string                 `json:"uuid"`
	UpdatedAt        int64                  `json:"updatedAt"`
	CreatedAt        int64                  `json:"createdAt"`
	ScannedAt        int64                  `json:"scannedAt"`
	Content          bool                   `json:"content"`
	Directory        bool                   `json:"directory"`
	ContentChangedAt int64                  `json:"contentChangedAt"`
	Hidden           int64                  `json:"hidden"`
	Locations        []*PlexLibraryLocation `json:"Location"`
}

type PlexLibraries struct {
	MediaContainer struct {
		Size        int64          `json:"size"`
		Directories []*PlexLibrary `json:"Directory"`
	} `json:"MediaContainer"`
}

func (p *Client) GetLibraries() ([]*PlexLibrary, error) {
	dest := &PlexLibraries{}
	resp, err := p._server().SetResult(dest).SetFormDataFromValues(p.data).Get("/library/sections")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, errors.Errorf("failed to get libraries: %s", resp.Status())
	}

	return dest.MediaContainer.Directories, nil
}

func (p *Client) LibraryType(section string) (int, error) {
	t, err := p.LibraryTypeName(section)
	if err != nil {
		return 0, err
	}
	if t == "" {
		return 0, errors.Errorf("library section %s not found", section)
	}

	id := p.LibraryTypeID(t)
	if id == PlexLibraryTypeUnknown {
		return 0, errors.Errorf("library section %s unknown", section)
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
	return "", errors.Errorf("library section %s not found", section)
}

func (p *Client) LibraryTypeID(t string) int {
	switch t {
	case "movie":
		return PlexLibraryTypeMovie
	case "show":
		return PlexLibraryTypeShow
	default:
		return PlexLibraryTypeUnknown
	}
}
