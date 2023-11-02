package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/dashotv/minion"
	"github.com/pkg/errors"
)

var workers *minion.Minion

type Job struct {
	Function minion.Func
	Schedule string
}

var jobs = map[string]Job{
	"PopularReleases": {PopularReleases, "0 */5 * * * *"},
	"CleanPlexPins":   {CleanPlexPins, "0 0 11 * * *"},
	"CleanJobs":       {CleanJobs, "0 0 11 * * *"},
	"PlexUserUpdates": {PlexUserUpdates, ""}, // "0 0 * * * *"
	// "DownloadsProcess": {DownloadsProcess, "*/5 * * * * *"},
}

func setupWorkers() error {
	workers = minion.New(cfg.Minion.Concurrency).WithLogger(log.Named("minion"))

	for n, j := range jobs {
		workers.Register(n, wrapJob(n, j.Function))
		if cfg.Cron {
			if j.Schedule != "" {
				if _, err := workers.Schedule(j.Schedule, n); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func CleanPlexPins() error {
	list, err := db.Pin.Query().
		GreaterThan("created_at", time.Now().UTC().AddDate(0, 0, -1)).
		Run()
	if err != nil {
		return errors.Wrap(err, "querying pins")
	}

	for _, p := range list {
		err := db.Pin.Delete(p)
		if err != nil {
			return errors.Wrap(err, "deleting pin")
		}
	}

	return nil
}

func CleanJobs() error {
	list, err := db.MinionJob.Query().
		GreaterThan("created_at", time.Now().UTC().AddDate(0, 0, -1)).
		Run()
	if err != nil {
		return errors.Wrap(err, "querying jobs")
	}

	for _, j := range list {
		err := db.MinionJob.Delete(j)
		if err != nil {
			return errors.Wrap(err, "deleting job")
		}
	}

	return nil
}

func PopularReleases() error {
	limit := 25
	intervals := map[string]int{
		"daily":   1,
		"weekly":  7,
		"monthly": 30,
	}

	for f, i := range intervals {
		for _, t := range releaseTypes {
			date := time.Now().AddDate(0, 0, -i)

			results, err := db.ReleasesPopularQuery(t, date, limit)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("popular releases %s %s", f, t))
			}

			cache.Set(fmt.Sprintf("releases_popular_%s_%s", f, t), results)
		}
	}

	return nil
}

func PlexUserUpdates() error {
	server.Log.Debugf("PlexUserUpdates: updating users")

	pins, err := db.Pin.Query().Run()
	if err != nil {
		return errors.Wrap(err, "querying pins")
	}

	check := map[string]bool{}
	users := []*User{}
	server.Log.Debugf("PlexUserUpdates: ranging pins")
	for _, p := range pins {
		if p.Token == "" {
			continue
		}

		if !check[p.Token] {
			check[p.Token] = true
			server.Log.Debugf("PlexUserUpdates: find user by token %s", p.Token)
			resp, err := db.User.Query().Where("token", p.Token).Run()
			if err != nil {
				return errors.Wrap(err, "querying user")
			}
			if len(resp) > 1 {
				return errors.New("multiple users found")
			}

			if len(resp) == 1 {
				server.Log.Debugf("PlexUserUpdates: adding user")
				users = append(users, resp[0])
				continue
			}

			user := &User{Token: p.Token}
			err = db.User.Save(user)
			if err != nil {
				return errors.Wrap(err, "saving user")
			}
			server.Log.Debugf("PlexUserUpdates: adding new user")
			users = append(users, user)
		}
	}

	server.Log.Debugf("PlexUserUpdates: ranging users %d", len(users))
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

		server.Log.Debugf("PlexUserUpdates: updating user %s", u.Name)
		err = db.User.Update(u)
		if err != nil {
			return errors.Wrap(err, "updating user")
		}

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
			if d == nil || len(d.MediaContainer.Metadata) == 0 {
				continue
			}
			for _, dm := range d.MediaContainer.Metadata {
				m, err := findMediaByGUIDs(dm.GUID)
				if err != nil {
					return errors.Wrap(err, "finding media")
				}
				if m != nil {
					continue
				}
				server.Log.Debugf("PlexUserUpdates: %s: NOT FOUND", dm.Title)
				// TODO: create request for media
			}
		}
	}
	// upate each users watchlist data
	return nil
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

func CausingErrors() error {
	log.Info("causing error")
	return nil
}

func DownloadsProcess() error {
	log.Info("processing downloads")
	return nil
}

func ProcessFeeds() error {
	log.Info("processing feeds")
	return db.ProcessFeeds()
}

func wrapJob(name string, f func() error) func() error {
	return func() error {
		j := &MinionJob{Name: name}

		err := db.MinionJob.Save(j)
		if err != nil {
			return errors.Wrap(err, "saving job")
		}

		start := time.Now()
		ferr := f()
		if ferr != nil {
			log.Errorf("job:%s: %s", name, ferr)
			j.Error = ferr.Error()
		}

		j.ProcessedAt = time.Now()
		duration := time.Since(start)
		j.Duration = duration.Seconds()

		err = db.MinionJob.Update(j)
		if err != nil {
			return errors.Wrap(err, "updating job")
		}

		log.Infof("job:%s: %s", name, duration)
		return ferr
	}
}
