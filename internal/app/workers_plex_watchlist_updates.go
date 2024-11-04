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
		err = createRequest("rss", item.Title, item.Categories[0], item.GUID)
		if err != nil {
			// a.Log.Debugf("PlexUserUpdates: NOT FOUND: %s: %s %+v", dm.Title, dm.Type, dm.GUID)
			// notifier.Log.Warnf("Watchlist", "NOT FOUND: %s: %s %+v", dm.Title, dm.Type, dm.GUID)
			continue
		}
	}

	// 	users, err := a.DB.User.Query().NotEqual("token", "").Run()
	// 	if err != nil {
	// 		return fae.Wrap(err, "querying users")
	// 	}
	//
	// 	for _, u := range users {
	// 		list, err := a.Plex.GetWatchlist(u.Token)
	// 		if err != nil {
	// 			notifier.Log.Errorf("PlexWatchlistUpdates", "getting watchlist: %s: %s", u.Name, err)
	// 			continue
	// 		}
	//
	// 		if list == nil || len(list.MediaContainer.Metadata) == 0 {
	// 			continue
	// 		}
	//
	// 		details, err := a.Plex.GetWatchlistDetail(u.Token, list)
	// 		if err != nil {
	// 			return fae.Wrap(err, "getting watchlist details")
	// 		}
	//
	// 		for _, d := range details {
	// 			if d == nil || d.MediaContainer.Size != 1 {
	// 				a.Log.Debugf("PlexUserUpdates: dm empty? size %d len %d", d.MediaContainer.Size, len(d.MediaContainer.Metadata))
	// 				continue
	// 			}
	// 			dm := d.MediaContainer.Metadata[0]
	// 			m, err := findMediaByGUIDs(dm.GUID)
	// 			if err != nil {
	// 				return fae.Wrap(err, "finding media")
	// 			}
	// 			if m != nil {
	// 				continue
	// 			}
	// 			err = createRequest(u.Name, dm.Title, dm.Type, dm.GUID)
	// 			if err != nil {
	// 				// a.Log.Debugf("PlexUserUpdates: NOT FOUND: %s: %s %+v", dm.Title, dm.Type, dm.GUID)
	// 				// notifier.Log.Warnf("Watchlist", "NOT FOUND: %s: %s %+v", dm.Title, dm.Type, dm.GUID)
	// 				continue
	// 			}
	// 			// a.Log.Infof("PlexUserUpdates: REQUESTED: %s: %s", dm.Title, dm.Type)
	// 		}
	// 	}

	return a.Workers.Enqueue(&CreateMediaFromRequests{})
}

func createRequest(user, title, t string, guid string) error {
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
	q := app.DB.Request.Query()

	if source == "tmdb" {
		q = q.Where("source", "tmdb").Where("source_id", source_id)
	} else if source == "imdb" {
		q = q.Where("imdb_id", source_id)
	} else {
		return fae.Errorf("createMovieRequest: unknown source: %s", source)
	}

	reqs, err := q.Run()
	if err != nil {
		return fae.Wrap(err, "querying requests")
	}
	if len(reqs) > 0 {
		return nil
	}

	req := &Request{
		User:     user,
		Title:    title,
		Source:   source,
		SourceID: source_id,
		Type:     "movie",
		Status:   requestDefaultStatus,
	}

	err = app.DB.Request.Save(req)
	if err != nil {
		return fae.Wrap(err, "saving request")
	}

	return nil
}

func createShowRequest(user, title string, guid string) error {
	source_id := guidToSourceID("tvdb", guid)
	if source_id == "" {
		return fae.New("createShowRequest: no tvdb guid")
	}

	reqs, err := app.DB.Request.Query().Where("source", "tvdb").Where("source_id", source_id).Run()
	if err != nil {
		return fae.Wrap(err, "querying requests")
	}
	if len(reqs) > 0 {
		return nil
	}

	req := &Request{
		User:     user,
		Title:    title,
		Source:   "tvdb",
		SourceID: source_id,
		Type:     "series",
		Status:   requestDefaultStatus,
	}

	err = app.DB.Request.Save(req)
	if err != nil {
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
	s := strings.Split(guid, "://")
	list, err := app.DB.Medium.Query().Where("source", s[0]).Where("source_id", s[1]).Run()
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
