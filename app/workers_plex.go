package app

import (
	"strings"

	"github.com/pkg/errors"
)

// PlexPinToUsers ensures users are created from athorized pins
func PlexPinToUsers() error {
	log := log.Named("job.PlexPinToUsers")
	log.Debugf("creating users from authenticated pins")

	pins, err := db.Pin.Query().Run()
	if err != nil {
		return errors.Wrap(err, "querying pins")
	}

	check := map[string]bool{}
	log.Debugf("ranging pins")
	for _, p := range pins {
		if p.Token == "" {
			continue
		}

		if check[p.Token] {
			continue
		}

		check[p.Token] = true
		log.Debugf("find user by token %s", p.Token)
		resp, err := db.User.Query().Where("token", p.Token).Run()
		if err != nil {
			return errors.Wrap(err, "querying user")
		}
		if len(resp) > 0 {
			// users exists
			continue
		}

		// create user
		user := &User{Token: p.Token}
		err = db.User.Save(user)
		if err != nil {
			return errors.Wrap(err, "saving user")
		}
	}

	err = workers.Enqueue("PlexUserUpdates")
	if err != nil {
		return errors.Wrap(err, "enqueuing worker")
	}

	return nil
}

// PlexUserUpdates updates users from plex
func PlexUserUpdates() error {
	log := log.Named("job.PlexUserUpdates")
	log.Debugf("updating users")

	users, err := db.User.Query().NotEqual("token", "").Run()
	if err != nil {
		return errors.Wrap(err, "querying users")
	}

	for _, u := range users {
		data, err := plex.GetUser(u.Token)
		if err != nil {
			return errors.Wrap(err, "getting user data")
		}

		u.Name = data.Username
		u.Email = data.Email
		u.Thumb = data.Thumb
		u.Home = data.Home
		u.Admin = data.HomeAdmin

		log.Debugf("updating user %s", u.Name)
		err = db.User.Update(u)
		if err != nil {
			return errors.Wrap(err, "updating user")
		}
	}

	err = workers.Enqueue("PlexWatchlistUpdates")
	if err != nil {
		return errors.Wrap(err, "enqueuing worker")
	}

	return nil
}

// PlexWatchlistUpdate updates watchlist from plex
func PlexWatchlistUpdates() error {
	log := log.Named("job.PlexWatchlistUpdates")
	log.Debugf("creating requests from watchlists")

	users, err := db.User.Query().NotEqual("token", "").Run()
	if err != nil {
		return errors.Wrap(err, "querying users")
	}

	for _, u := range users {
		list, err := plex.GetWatchlist(u.Token)
		if err != nil {
			return errors.Wrap(err, "getting watchlist")
		}

		if list == nil || len(list.MediaContainer.Metadata) == 0 {
			continue
		}

		details, err := plex.GetWatchlistDetail(u.Token, list)
		if err != nil {
			return errors.Wrap(err, "getting watchlist details")
		}

		for _, d := range details {
			if d == nil || d.MediaContainer.Size != 1 {
				log.Debugf("PlexUserUpdates: dm empty? size %d len %d", d.MediaContainer.Size, len(d.MediaContainer.Metadata))
				continue
			}
			dm := d.MediaContainer.Metadata[0]
			m, err := findMediaByGUIDs(dm.GUID)
			if err != nil {
				return errors.Wrap(err, "finding media")
			}
			if m != nil {
				continue
			}
			log.Debugf("PlexUserUpdates: NOT FOUND: %s: %s", dm.Title, dm.Type)
			err = createRequest(u.Name, dm.Title, dm.Type, dm.GUID)
			if err != nil {
				return errors.Wrap(err, "creating request")
			}
			log.Infof("PlexUserUpdates: REQUESTED: %s: %s", dm.Title, dm.Type)
		}
	}
	return nil
}

func createRequest(user, title, t string, guids []GUID) error {
	switch t {
	case "movie":
		return createMovieRequest(user, title, guids)
	case "show":
		return createShowRequest(user, title, guids)
	default:
		return errors.Errorf("createRequest: unknown type: %s", t)
	}
}

func createMovieRequest(user, title string, guids []GUID) error {
	source_id := guidToSourceID("tmdb", guids)
	if source_id == "" {
		return errors.New("createMovieRequest: no tmdb guid")
	}

	reqs, err := db.Request.Query().Where("source", "tmdb").Where("source_id", source_id).Run()
	if err != nil {
		return errors.Wrap(err, "querying requests")
	}
	if len(reqs) > 0 {
		return nil
	}

	req := &Request{
		User:     user,
		Title:    title,
		Source:   "tmdb",
		SourceId: source_id,
		Type:     "movie",
	}

	err = db.Request.Save(req)
	if err != nil {
		return errors.Wrap(err, "saving request")
	}

	return nil
}

func createShowRequest(user, title string, guids []GUID) error {
	source_id := guidToSourceID("tvdb", guids)
	if source_id == "" {
		return errors.New("createShowRequest: no tvdb guid")
	}

	reqs, err := db.Request.Query().Where("source", "tvdb").Where("source_id", source_id).Run()
	if err != nil {
		return errors.Wrap(err, "querying requests")
	}
	if len(reqs) > 0 {
		return nil
	}

	req := &Request{
		User:     user,
		Title:    title,
		Source:   "tvdb",
		SourceId: source_id,
		Type:     "series",
	}

	err = db.Request.Save(req)
	if err != nil {
		return errors.Wrap(err, "saving request")
	}
	return nil
}

func guidToSourceID(source string, guids []GUID) string {
	for _, g := range guids {
		s := strings.Split(g.ID, "://")
		if s[0] == source {
			return s[1]
		}
	}

	return ""
}

func findMediaByGUIDs(list []GUID) (*Medium, error) {
	for _, g := range list {
		s := strings.Split(g.ID, "://")
		list, err := db.Medium.Query().Where("source", s[0]).Where("source_id", s[1]).Run()
		if err != nil {
			return nil, errors.Wrap(err, "querying media")
		}
		if len(list) > 0 {
			return list[0], nil
		}
	}

	return nil, nil
}
