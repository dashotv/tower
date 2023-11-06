package app

import (
	"fmt"
	"time"

	"github.com/dashotv/minion"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

var workers *minion.Minion

type Job struct {
	Function minion.Func
	Schedule string
}

var jobs = map[string]Job{
	"PopularReleases":          {PopularReleases, "0 */5 * * * *"},       // every 5 minutes
	"CleanPlexPins":            {CleanPlexPins, "0 0 11 * * *"},          // every day at 11am UTC
	"CleanJobs":                {CleanJobs, "0 0 11 * * *"},              // every day at 11am UTC
	"CleanLogs":                {CleanLogs, "0 0 11 * * *"},              // every day at 11am UTC
	"PlexPinToUsers":           {PlexPinToUsers, ""},                     // run on demand from route
	"PlexUserUpdates":          {PlexUserUpdates, "0 0 11 * * *"},        // every day at 11am UTC
	"PlexWatchlistUpdates":     {PlexWatchlistUpdates, "0 0 * * * *"},    // every hour and on demond
	"CreateMediaFromRequests":  {CreateMediaFromRequests, "0 0 * * * *"}, // every hour and on demand
	"TmdbUpdateMovie":          {TmdbUpdateMovie, ""},                    // run on demoand
	"TmdbUpdateMovieImage":     {TmdbUpdateMovieImage, ""},               // run on demand
	"TvdbUpdateSeries":         {TvdbUpdateSeries, ""},                   // run on demand
	"TvdbUpdateSeriesImage":    {TvdbUpdateSeriesImage, ""},              // run on demand
	"TvdbUpdateSeriesEpisodes": {TvdbUpdateSeriesEpisodes, ""},           // run on demand
	// "DownloadsProcess": {DownloadsProcess, "*/5 * * * * *"},
}

func setupWorkers() error {
	workers = minion.New(cfg.Minion.Concurrency).WithLogger(log.Named("minion"))

	for n, j := range jobs {
		if cfg.Cron && j.Schedule != "" {
			if _, err := workers.Schedule(j.Schedule, n, wrapJob(n, j.Function)); err != nil {
				return err
			}
		} else {
			workers.Register(n, wrapJob(n, j.Function))
		}
	}

	return nil
}

func CleanPlexPins(_ any) error {
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

func CleanJobs(_ any) error {
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

func CleanLogs(_ any) error {
	// _, err := db.Message.Collection.DeleteMany(context.Background(), bson.M{"created_at": bson.M{"$lt": time.Now().UTC().AddDate(0, 0, -1)}})
	// if err != nil {
	// 	return errors.Wrap(err, "cleaning logs")
	// }

	return nil
}

func PopularReleases(_ any) error {
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

func CreateMediaFromRequests(_ any) error {
	log := log.Named("job.CreateMediaFromRequests")

	requests, err := db.Request.Query().Where("status", "approved").Run()
	if err != nil {
		return errors.Wrap(err, "querying requests")
	}

	for _, r := range requests {
		log.Infof("processing request: %s", r.Title)
		if r.Source == "tmdb" {
			err := createMovieFromRequest(r)
			if err != nil {
				log.Errorf("creating movie from request: %s", err)
				r.Status = "failed"
			} else {
				log.Infof("created movie: %s", r.Title)
				r.Status = "completed"
			}
		} else if r.Source == "tvdb" {
			err := createShowFromRequest(r)
			if err != nil {
				log.Errorf("creating series from request: %s", err)
				r.Status = "failed"
			} else {
				log.Infof("created series: %s", r.Title)
				r.Status = "completed"
			}
		}

		log.Infof("request: [%s] %s", r.Status, r.Title)
		if err := db.Request.Update(r); err != nil {
			return errors.Wrap(err, "updating request")
		}

		if err := events.Send("tower.requests", &EventTowerRequest{Event: "update", ID: r.ID.Hex(), Request: r}); err != nil {
			return errors.Wrap(err, "sending event")
		}
	}
	return nil
}

func createShowFromRequest(r *Request) error {
	count, err := db.Series.Count(bson.M{"_type": "Series", "source": r.Source, "source_id": r.SourceId})
	if err != nil {
		return errors.Wrap(err, "counting series")
	}
	if count > 0 {
		return nil
	}

	s := &Series{
		Type:     "Series",
		Source:   r.Source,
		SourceId: r.SourceId,
		Title:    r.Title,
		Kind:     "tv",
	}

	err = db.Series.Save(s)
	if err != nil {
		return errors.Wrap(err, "saving show")
	}

	if err := workers.EnqueueWithPayload("TvdbUpdateSeries", s.ID.Hex()); err != nil {
		return errors.Wrap(err, "queueing update job")
	}
	return nil
}

func createMovieFromRequest(r *Request) error {
	count, err := db.Series.Count(bson.M{"_type": "Movie", "source": r.Source, "source_id": r.SourceId})
	if err != nil {
		return errors.Wrap(err, "counting series")
	}
	if count > 0 {
		return nil
	}

	m := &Movie{
		Type:     "Movie",
		Source:   r.Source,
		SourceId: r.SourceId,
		Title:    r.Title,
		Kind:     "movies",
	}

	err = db.Movie.Save(m)
	if err != nil {
		return errors.Wrap(err, "saving movie")
	}

	if err := workers.EnqueueWithPayload("TmdbUpdateMovie", m.ID.Hex()); err != nil {
		return errors.Wrap(err, "queueing update job")
	}
	return nil
}

func CausingErrors(_ any) error {
	log.Info("causing error")
	return nil
}

func DownloadsProcess(_ any) error {
	log.Info("processing downloads")
	return nil
}

func ProcessFeeds(_ any) error {
	log.Info("processing feeds")
	return db.ProcessFeeds()
}

func wrapJob(name string, f minion.Func) minion.Func {
	return func(payload any) error {
		j := &MinionJob{Name: name}

		err := db.MinionJob.Save(j)
		if err != nil {
			return errors.Wrap(err, "saving job")
		}

		start := time.Now()
		ferr := f(payload)
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
