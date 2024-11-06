package app

import (
	"context"
	"strings"

	"github.com/mmcdole/gofeed"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
	"github.com/dashotv/tower/internal/plex"
)

var requestDefaultStatus = "approved"

// PlexWatchlistUpdates updates watchlist from plex
type PlexWatchlistUpdates struct {
	minion.WorkerDefaults[*PlexWatchlistUpdates]
}

func (j *PlexWatchlistUpdates) Kind() string { return "PlexWatchlistUpdates" }
func (j *PlexWatchlistUpdates) Work(ctx context.Context, job *minion.Job[*PlexWatchlistUpdates]) error {
	a := ContextApp(ctx)
	if a == nil {
		return fae.New("PlexWatchlistUpdates: no app in context")
	}

	url := a.Config.PlexWatchlistURL
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return fae.Wrap(err, "parsing feed")
	}

	for _, item := range feed.Items {
		a.Log.Debugf("PlexWatchlistUpdates: %s %+v %s", item.Title, item.Categories, item.GUID)
		m, err := findMediaByGuid(item.GUID)
		if err != nil {
			return fae.Wrap(err, "finding media")
		}
		if m != nil {
			continue
		}
		if err := createRequest("rss", item.Title, item.Categories[0], item.GUID); err != nil {
			a.Log.Debugf("PlexUserUpdates: NOT FOUND: %s: %s %+v", item.Title, item.Categories[0], item.GUID)
			// notifier.Log.Warnf("Watchlist", "NOT FOUND: %s: %s %+v", dm.Title, dm.Type, dm.GUID)
			continue
		}
	}

	return a.Workers.Enqueue(&CreateMediaFromRequests{})
}

func createRequest(user, title, t string, guid string) error {
	exists, err := app.DB.RequestExists(guid)
	if err != nil {
		return fae.Wrap(err, "checking request exists")
	}
	if exists {
		return nil
	}

	switch t {
	case "movie":
		return createMovieRequest(user, title, guid)
	case "show":
		return createShowRequest(user, title, guid)
	default:
		return fae.Errorf("createRequest: unknown type: %s", t)
	}
}

func createMovieRequest(user, title string, guid string) error {
	source, source_id := guidSplit(guid)

	req := &Request{
		User:     user,
		Title:    title,
		Source:   source,
		SourceID: source_id,
		Type:     "movie",
		Status:   requestDefaultStatus,
	}

	if err := app.DB.Request.Save(req); err != nil {
		return fae.Wrap(err, "saving request")
	}

	return nil
}

func createShowRequest(user, title string, guid string) error {
	source, source_id := guidSplit(guid)

	req := &Request{
		User:     user,
		Title:    title,
		Source:   source,
		SourceID: source_id,
		Type:     "series",
		Status:   requestDefaultStatus,
	}

	if err := app.DB.Request.Save(req); err != nil {
		return fae.Wrap(err, "saving request")
	}

	return nil
}

func guidSplit(guid string) (string, string) {
	s := strings.Split(guid, "://")
	if len(s) != 2 {
		return "", ""
	}

	return s[0], s[1]
}

func guidToSourceID(source string, guid string) string {
	s := strings.Split(guid, "://")
	if s[0] == source {
		return s[1]
	}

	return ""
}

func findMediaByGuid(guid string) (*Medium, error) {
	source, source_id := guidSplit(guid)
	q := app.DB.Medium.Query()

	if source == "tmdb" || source == "tvdb" {
		q = q.Where("source", source).Where("source_id", source_id)
	} else if source == "imdb" {
		q = q.Where("imdb_id", source_id)
	} else {
		return nil, fae.Errorf("findMediaByGuid: unknown source: %s", source)
	}

	list, err := q.Run()
	if err != nil {
		return nil, fae.Wrap(err, "querying media")
	}
	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func findMediaByGUIDs(list []plex.GUID) (*Medium, error) {
	for _, g := range list {
		s := strings.Split(g.ID, "://")
		list, err := app.DB.Medium.Query().Where("source", s[0]).Where("source_id", s[1]).Run()
		if err != nil {
			return nil, fae.Wrap(err, "querying media")
		}
		if len(list) > 0 {
			return list[0], nil
		}
	}

	return nil, nil
}
