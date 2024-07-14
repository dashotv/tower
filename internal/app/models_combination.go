package app

import (
	"fmt"
	"sort"

	"github.com/dashotv/fae"
)

func (c *Connector) CombinationGet(id string) (*Combination, error) {
	m, err := c.Combination.Get(id, &Combination{})
	if err != nil {
		return nil, err
	}

	// post process here

	return m, nil
}

func (c *Connector) CombinationByName(name string) (*Combination, error) {
	list, err := c.Combination.Query().Where("name", name).Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, err
	}
	if len(list) > 1 {
		return nil, err
	}
	return list[0], nil
}

func (c *Connector) CombinationList(page, limit int) ([]*Combination, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	list, err := c.Combination.Query().Skip((page - 1) * limit).Limit(limit).Run()
	if err != nil {
		return nil, err
	}

	return list, nil
}

// type CombinationChild struct {
// 	RatingKey    string `json:"ratingKey"`
// 	Key          string `json:"key"`
// 	GUID         string `json:"guid"`
// 	Type         string `json:"type"`
// 	Title        string `json:"title"`
// 	LibraryID    int64  `json:"librarySectionID"`
// 	LibraryTitle string `json:"librarySectionTitle"`
// 	LibraryKey   string `json:"librarySectionKey"`
// 	Summary      string `json:"summary"`
// 	Thumb        string `json:"thumb"`
// 	Total        int    `json:"total"`
// 	Viewed       int    `json:"viewed"`
// 	Link         string `json:"link"`
// 	Next         string `json:"next"`
// 	LastViewedAt int64  `json:"lastViewedAt"`
// 	AddedAt      int64  `json:"addedAt"`
// 	UpdatedAt    int64  `json:"updatedAt"`
// }

type stuffSorter struct {
	list []*CombinationChild
	by   func(p1, p2 *CombinationChild) bool
}

// Len is part of sort.Interface.
func (s *stuffSorter) Len() int {
	return len(s.list)
}

// Swap is part of sort.Interface.
func (s *stuffSorter) Swap(i, j int) {
	s.list[i], s.list[j] = s.list[j], s.list[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *stuffSorter) Less(i, j int) bool {
	return s.by(s.list[i], s.list[j])
}

func (c *Connector) CombinationChildren(name string) ([]*CombinationChild, error) {
	comb, err := c.CombinationByName(name)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, id := range comb.Collections {
		col, err := c.CollectionGet(id)
		if err != nil {
			return nil, err
		}
		keys = append(keys, col.RatingKey)
	}

	list := []*CombinationChild{}

	for _, i := range keys {
		children, err := app.Plex.GetCollectionChildren(i)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			metadata, err := app.Plex.GetViewedByKey(child.RatingKey)
			if err != nil {
				return nil, err
			}
			if metadata == nil {
				return nil, fae.Errorf("metadata not found for %s", child.RatingKey)
			}
			if metadata.Leaves > 0 && metadata.Leaves == metadata.Viewed {
				continue
			}

			stuff := &CombinationChild{
				RatingKey:    child.RatingKey,
				Key:          child.Key,
				GUID:         child.GUID,
				Type:         child.Type,
				Title:        child.Title,
				LibraryID:    child.LibraryID,
				LibraryTitle: child.LibraryTitle,
				LibraryKey:   child.LibraryKey,
				Summary:      child.Summary,
				Link:         fmt.Sprintf("https://app.plex.tv/desktop/#!/server/%s/details?key=%s?X-Plex-Token=%s", app.Config.PlexMachineIdentifier, child.Key, app.Config.PlexToken),
				Thumb:        fmt.Sprintf("%s%s?X-Plex-Token=%s", app.Config.PlexServerURL, child.Thumb, app.Config.PlexToken),
				Total:        metadata.Leaves,
				Viewed:       metadata.Viewed,
				LastViewedAt: metadata.LastViewedAt,
			}

			if child.Type == "show" {
				unwatched, err := app.Plex.GetSeriesEpisodesUnwatched(child.RatingKey)
				if err != nil {
					return nil, err
				}
				if unwatched != nil {
					stuff.Next = unwatched.RatingKey
					stuff.AddedAt = unwatched.AddedAt
					stuff.UpdatedAt = unwatched.UpdatedAt
				}
			} else {
				stuff.AddedAt = child.AddedAt
				stuff.UpdatedAt = child.UpdatedAt
			}
			list = append(list, stuff)
		}
	}

	sorter := &stuffSorter{
		list: list,
		by: func(p1, p2 *CombinationChild) bool {
			return p1.AddedAt > p2.AddedAt
		},
	}
	sort.Sort(sorter)

	return list, nil
}
