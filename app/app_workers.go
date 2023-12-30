// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"context"

	"github.com/dashotv/minion"
	"github.com/pkg/errors"
)

func init() {
	initializers = append(initializers, setupWorkers)
	healthchecks["workers"] = checkWorkers
}

func checkWorkers(app *Application) error {
	// TODO: workers health check
	return nil
}

func setupWorkers(app *Application) error {
	ctx := context.Background()

	mcfg := &minion.Config{
		Logger:      app.Log.Named("minion"),
		Debug:       app.Config.MinionDebug,
		Concurrency: app.Config.MinionConcurrency,
		BufferSize:  app.Config.MinionBufferSize,
		DatabaseURI: app.Config.MinionURI,
		Database:    app.Config.MinionDatabase,
		Collection:  app.Config.MinionCollection,
	}

	m, err := minion.New(ctx, mcfg)
	if err != nil {
		return errors.Wrap(err, "creating minion")
	}

	// add something like the below line in app.Start() (before the workers are
	// started) to subscribe to job notifications.
	// minion sends notifications as jobs are processed and change status
	// m.Subscribe(app.MinionNotification)
	// an example of the subscription function and the basic setup instructions
	// are included at the end of this file.

	m.Queue("paths", 3, 3, 0)

	if err := minion.Register[*CleanPlexPins](m, &CleanPlexPins{}); err != nil {
		return errors.Wrap(err, "registering worker: clean_plex_pins (CleanPlexPins)")
	}
	if _, err := m.Schedule("0 0 11 * * *", &CleanPlexPins{}); err != nil {
		return errors.Wrap(err, "scheduling worker: clean_plex_pins (CleanPlexPins)")
	}

	if err := minion.Register[*CleanupJobs](m, &CleanupJobs{}); err != nil {
		return errors.Wrap(err, "registering worker: cleanup_jobs (CleanupJobs)")
	}
	if _, err := m.Schedule("0 10 11 * * *", &CleanupJobs{}); err != nil {
		return errors.Wrap(err, "scheduling worker: cleanup_jobs (CleanupJobs)")
	}

	if err := minion.Register[*CleanupLogs](m, &CleanupLogs{}); err != nil {
		return errors.Wrap(err, "registering worker: cleanup_logs (CleanupLogs)")
	}
	if _, err := m.Schedule("0 20 11 * * *", &CleanupLogs{}); err != nil {
		return errors.Wrap(err, "scheduling worker: cleanup_logs (CleanupLogs)")
	}

	if err := minion.Register[*CreateMediaFromRequests](m, &CreateMediaFromRequests{}); err != nil {
		return errors.Wrap(err, "registering worker: create_media_from_requests (CreateMediaFromRequests)")
	}
	if _, err := m.Schedule("15 0 * * * *", &CreateMediaFromRequests{}); err != nil {
		return errors.Wrap(err, "scheduling worker: create_media_from_requests (CreateMediaFromRequests)")
	}

	if err := minion.Register[*DownloadsProcess](m, &DownloadsProcess{}); err != nil {
		return errors.Wrap(err, "registering worker: downloads_process (DownloadsProcess)")
	}

	if err := minion.Register[*MediaPaths](m, &MediaPaths{}); err != nil {
		return errors.Wrap(err, "registering worker: media_paths (MediaPaths)")
	}

	if err := minion.RegisterWithQueue[*PathImport](m, &PathImport{}, "paths"); err != nil {
		return errors.Wrap(err, "registering worker: path_import (PathImport)")
	}

	if err := minion.Register[*PlexPinToUsers](m, &PlexPinToUsers{}); err != nil {
		return errors.Wrap(err, "registering worker: plex_pin_to_users (PlexPinToUsers)")
	}

	if err := minion.Register[*PlexUserUpdates](m, &PlexUserUpdates{}); err != nil {
		return errors.Wrap(err, "registering worker: plex_user_updates (PlexUserUpdates)")
	}
	if _, err := m.Schedule("0 30 11 * * *", &PlexUserUpdates{}); err != nil {
		return errors.Wrap(err, "scheduling worker: plex_user_updates (PlexUserUpdates)")
	}

	if err := minion.Register[*PlexWatchlistUpdates](m, &PlexWatchlistUpdates{}); err != nil {
		return errors.Wrap(err, "registering worker: plex_watchlist_updates (PlexWatchlistUpdates)")
	}
	if _, err := m.Schedule("0 0 * * * *", &PlexWatchlistUpdates{}); err != nil {
		return errors.Wrap(err, "scheduling worker: plex_watchlist_updates (PlexWatchlistUpdates)")
	}

	if err := minion.Register[*PopularReleases](m, &PopularReleases{}); err != nil {
		return errors.Wrap(err, "registering worker: popular_releases (PopularReleases)")
	}
	if _, err := m.Schedule("0 */5 * * * *", &PopularReleases{}); err != nil {
		return errors.Wrap(err, "scheduling worker: popular_releases (PopularReleases)")
	}

	if err := minion.Register[*TmdbUpdateAll](m, &TmdbUpdateAll{}); err != nil {
		return errors.Wrap(err, "registering worker: tmdb_update_all (TmdbUpdateAll)")
	}

	if err := minion.Register[*TmdbUpdateMovie](m, &TmdbUpdateMovie{}); err != nil {
		return errors.Wrap(err, "registering worker: tmdb_update_movie (TmdbUpdateMovie)")
	}

	if err := minion.Register[*TmdbUpdateMovieImage](m, &TmdbUpdateMovieImage{}); err != nil {
		return errors.Wrap(err, "registering worker: tmdb_update_movie_image (TmdbUpdateMovieImage)")
	}

	if err := minion.Register[*TvdbUpdateSeries](m, &TvdbUpdateSeries{}); err != nil {
		return errors.Wrap(err, "registering worker: tvdb_update_series (TvdbUpdateSeries)")
	}

	if err := minion.Register[*TvdbUpdateSeriesEpisodes](m, &TvdbUpdateSeriesEpisodes{}); err != nil {
		return errors.Wrap(err, "registering worker: tvdb_update_series_episodes (TvdbUpdateSeriesEpisodes)")
	}

	if err := minion.Register[*TvdbUpdateSeriesImage](m, &TvdbUpdateSeriesImage{}); err != nil {
		return errors.Wrap(err, "registering worker: tvdb_update_series_image (TvdbUpdateSeriesImage)")
	}

	if err := minion.Register[*UpdateIndexes](m, &UpdateIndexes{}); err != nil {
		return errors.Wrap(err, "registering worker: update_indexes (UpdateIndexes)")
	}

	app.Workers = m
	return nil
}

// run the following commands to create the events channel and add the necessary models.
//
// > golem add event jobs event id job:*Minion
// > golem add model minion_attempt --struct started_at:time.Time duration:float64 status error 'stacktrace:[]string'
// > golem add model minion queue kind args status 'attempts:[]*MinionAttempt'
//
// then add a Connection configuration that points to the same database connection information
// as the minion database.

// // This allows you to notify other services as jobs change status.
//func (a *Application) MinionNotification(n *minion.Notification) {
//	if n.JobID == "-" {
//		return
//	}
//
//	j := &Minion{}
//	err := app.DB.Minion.Find(n.JobID, j)
//	if err != nil {
//		log.Errorf("finding job: %s", err)
//		return
//	}
//
//	if n.Event == "job:created" {
//		events.Send("tower.jobs", &EventJob{"created", j.ID.Hex(), j})
//		return
//	}
//	events.Send("tower.jobs", &EventJob{"updated", j.ID.Hex(), j})
//}