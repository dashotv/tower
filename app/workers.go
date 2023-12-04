package app

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/dashotv/minion"
)

var workers *minion.Minion
var workersList = map[string]minion.Payload{
	"CleanupLogs":              &CleanupLogs{},
	"CleanupJobs":              &CleanupJobs{},
	"PopularReleases":          &PopularReleases{},
	"CleanPlexPins":            &CleanPlexPins{},
	"PlexPinToUsers":           &PlexPinToUsers{},
	"PlexUserUpdates":          &PlexUserUpdates{},
	"PlexWatchlistUpdates":     &PlexWatchlistUpdates{},
	"CreateMediaFromRequests":  &CreateMediaFromRequests{},
	"TmdbUpdateMovie":          &TmdbUpdateMovie{},
	"TmdbUpdateMovieImage":     &TmdbUpdateMovieImage{},
	"TvdbUpdateSeries":         &TvdbUpdateSeries{},
	"TvdbUpdateSeriesImage":    &TvdbUpdateSeriesImage{},
	"TvdbUpdateSeriesEpisodes": &TvdbUpdateSeriesEpisodes{},
	"DownloadsProcess":         &DownloadsProcess{},
	// "DownloadsFileMove":        &DownloadFileMover{},
}

func setupWorkers() error {
	ctx := context.Background()

	mcfg := &minion.Config{
		Debug:       true,
		Concurrency: cfg.Minion.Concurrency,
		Logger:      log.Named("minion"),
		Database:    cfg.Connections["minion"].Database,
		Collection:  cfg.Connections["minion"].Collection,
		DatabaseURI: cfg.Connections["minion"].URI,
	}
	m, err := minion.New(ctx, mcfg)
	if err != nil {
		return errors.Wrap(err, "creating minion")
	}

	m.Subscribe(func(n *minion.Notification) {
		if n.Event != "job:start" && n.Event != "job:success" && n.Event != "job:fail" {
			return
		}

		j := &Minion{}
		err := db.Minion.Find(n.JobID, j)
		if err != nil {
			log.Errorf("finding job: %s", err)
			return
		}

		if n.Event == "job:start" {
			events.Send("tower.jobs", &EventTowerJob{"created", j.ID.Hex(), j})
			return
		}
		events.Send("tower.jobs", &EventTowerJob{"updated", j.ID.Hex(), j})
	})

	if err := minion.Register[*CleanupLogs](m, &CleanupLogs{}); err != nil {
		return errors.Wrap(err, "registering worker")
	}
	if err := minion.Register[*CleanupJobs](m, &CleanupJobs{}); err != nil {
		return errors.Wrap(err, "registering worker: CleanupJobs")
	}
	if err := minion.Register[*PopularReleases](m, &PopularReleases{}); err != nil {
		return errors.Wrap(err, "registering worker: PopularReleases")
	}
	if err := minion.Register[*CleanPlexPins](m, &CleanPlexPins{}); err != nil {
		return errors.Wrap(err, "registering worker: CleanPlexPins")
	}
	if err := minion.Register[*PlexPinToUsers](m, &PlexPinToUsers{}); err != nil {
		return errors.Wrap(err, "registering worker: PlexPinToUsers")
	}
	if err := minion.Register[*PlexUserUpdates](m, &PlexUserUpdates{}); err != nil {
		return errors.Wrap(err, "registering worker: PlexUserUpdates")
	}
	if err := minion.Register[*PlexWatchlistUpdates](m, &PlexWatchlistUpdates{}); err != nil {
		return errors.Wrap(err, "registering worker: PlexWatchlistUpdates")
	}
	if err := minion.Register[*CreateMediaFromRequests](m, &CreateMediaFromRequests{}); err != nil {
		return errors.Wrap(err, "registering worker: CreateMediaFromRequests")
	}
	if err := minion.Register[*TmdbUpdateMovie](m, &TmdbUpdateMovie{}); err != nil {
		return errors.Wrap(err, "registering worker: TmdbUpdateMovie")
	}
	if err := minion.Register[*TmdbUpdateMovieImage](m, &TmdbUpdateMovieImage{}); err != nil {
		return errors.Wrap(err, "registering worker: TmdbUpdateMovieImage")
	}
	if err := minion.Register[*TvdbUpdateSeries](m, &TvdbUpdateSeries{}); err != nil {
		return errors.Wrap(err, "registering worker: TvdbUpdateSeries")
	}
	if err := minion.Register[*TvdbUpdateSeriesImage](m, &TvdbUpdateSeriesImage{}); err != nil {
		return errors.Wrap(err, "registering worker: TvdbUpdateSeriesImage")
	}
	if err := minion.Register[*TvdbUpdateSeriesEpisodes](m, &TvdbUpdateSeriesEpisodes{}); err != nil {
		return errors.Wrap(err, "registering worker: TvdbUpdateSeriesEpisodes")
	}
	if err := minion.Register[*DownloadsProcess](m, &DownloadsProcess{}); err != nil {
		return errors.Wrap(err, "registering worker: DownloadsProcess")
	}
	// if err := minion.Register[*DownloadFileMover](m, &DownloadFileMover{}); err != nil {
	// 	return errors.Wrap(err, "registering worker: DownloadsProcess")
	// }

	if _, err := m.Schedule("0 */5 * * * *", &PopularReleases{}); err != nil {
		return errors.Wrap(err, "scheduling worker: PopularReleases")
	}
	if _, err := m.Schedule("0 0 11 * * *", &CleanPlexPins{}); err != nil {
		return errors.Wrap(err, "scheduling worker: CleanPlexPins")
	}
	if _, err := m.Schedule("0 10 11 * * *", &CleanupJobs{}); err != nil {
		return errors.Wrap(err, "scheduling worker: CleanJobs")
	}
	if _, err := m.Schedule("0 20 11 * * *", &CleanupLogs{}); err != nil {
		return errors.Wrap(err, "scheduling worker: CleanLogs")
	}
	if _, err := m.Schedule("0 30 11 * * *", &PlexUserUpdates{}); err != nil {
		return errors.Wrap(err, "scheduling worker: PlexUserUpdates")
	}
	if _, err := m.Schedule("0 0 * * * *", &PlexWatchlistUpdates{}); err != nil {
		return errors.Wrap(err, "scheduling worker: PlexWatchlistUpdates")
	}
	if _, err := m.Schedule("15 0 * * * *", &CreateMediaFromRequests{}); err != nil {
		return errors.Wrap(err, "scheduling worker: CreateMediaFromRequests")
	}

	// if err := minion.Register[*Ping](m, &Ping{}); err != nil {
	// 	return errors.Wrap(err, "registering worker: Ping")
	// }
	// if _, err := m.Schedule("*/30 * * * * *", &Ping{}); err != nil {
	// 	return errors.Wrap(err, "scheduling worker: Ping")
	// }

	workers = m
	return nil
}

type Ping struct{}

func (j *Ping) Kind() string { return "ping" }
func (j *Ping) Work(ctx context.Context, job *minion.Job[*Ping]) error {
	log.Named("ping").Debug("ping")
	time.Sleep(8 * time.Second)
	// return errors.Errorf("testing error: %s", time.Now())
	return nil
}
