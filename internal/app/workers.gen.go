// Code generated by github.com/dashotv/golem. DO NOT EDIT.
package app

import (
	"context"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

func init() {
	initializers = append(initializers, setupWorkers)
	healthchecks["workers"] = checkWorkers
	starters = append(starters, startWorkers)
}

func checkWorkers(app *Application) error {
	// TODO: workers health check
	return nil
}

func startWorkers(ctx context.Context, app *Application) error {
	ctx = context.WithValue(ctx, "app", app)

	app.Log.Debugf("starting workers (%d)", app.Config.MinionConcurrency)
	go app.Workers.Start(ctx)

	return nil
}

func setupWorkers(app *Application) error {
	mcfg := &minion.Config{
		Logger:      app.Log.Named("minion"),
		Debug:       app.Config.MinionDebug,
		Concurrency: app.Config.MinionConcurrency,
		BufferSize:  app.Config.MinionBufferSize,
		DatabaseURI: app.Config.MinionURI,
		Database:    app.Config.MinionDatabase,
		Collection:  app.Config.MinionCollection,
	}

	m, err := minion.New("tower", mcfg)
	if err != nil {
		return fae.Wrap(err, "creating minion")
	}

	// add something like the below line in app.Start() (before the workers are
	// started) to subscribe to job notifications.
	// minion sends notifications as jobs are processed and change status
	// m.Subscribe(app.MinionNotification)
	// an example of the subscription function and the basic setup instructions
	// are included at the end of this file.

	m.Queue("paths", 3, 3, 0)
	m.Queue("series", 3, 0, 5)

	if err := minion.Register[*CleanPlexPins](m, &CleanPlexPins{}); err != nil {
		return fae.Wrap(err, "registering worker: clean_plex_pins (CleanPlexPins)")
	}
	if _, err := m.Schedule("0 0 11 * * *", &CleanPlexPins{}); err != nil {
		return fae.Wrap(err, "scheduling worker: clean_plex_pins (CleanPlexPins)")
	}

	if err := minion.Register[*CleanupLogs](m, &CleanupLogs{}); err != nil {
		return fae.Wrap(err, "registering worker: cleanup_logs (CleanupLogs)")
	}
	if _, err := m.Schedule("0 20 11 * * *", &CleanupLogs{}); err != nil {
		return fae.Wrap(err, "scheduling worker: cleanup_logs (CleanupLogs)")
	}

	if err := minion.Register[*CreateMediaFromRequests](m, &CreateMediaFromRequests{}); err != nil {
		return fae.Wrap(err, "registering worker: create_media_from_requests (CreateMediaFromRequests)")
	}
	if _, err := m.Schedule("15 0 * * * *", &CreateMediaFromRequests{}); err != nil {
		return fae.Wrap(err, "scheduling worker: create_media_from_requests (CreateMediaFromRequests)")
	}

	if err := minion.Register[*DownloadsMovies](m, &DownloadsMovies{}); err != nil {
		return fae.Wrap(err, "registering worker: downloads_movies (DownloadsMovies)")
	}
	if _, err := m.Schedule("0 0 * * * *", &DownloadsMovies{}); err != nil {
		return fae.Wrap(err, "scheduling worker: downloads_movies (DownloadsMovies)")
	}

	if err := minion.Register[*DownloadsProcess](m, &DownloadsProcess{}); err != nil {
		return fae.Wrap(err, "registering worker: downloads_process (DownloadsProcess)")
	}
	if _, err := m.Schedule("0 * * * * *", &DownloadsProcess{}); err != nil {
		return fae.Wrap(err, "scheduling worker: downloads_process (DownloadsProcess)")
	}

	if err := minion.Register[*DownloadsProcessLoad](m, &DownloadsProcessLoad{}); err != nil {
		return fae.Wrap(err, "registering worker: downloads_process_load (DownloadsProcessLoad)")
	}

	if err := minion.RegisterWithQueue[*FileWalk](m, &FileWalk{}, "paths"); err != nil {
		return fae.Wrap(err, "registering worker: file_walk (FileWalk)")
	}
	if _, err := m.Schedule("0 0 * * * *", &FileWalk{}); err != nil {
		return fae.Wrap(err, "scheduling worker: file_walk (FileWalk)")
	}

	if err := minion.Register[*FilesMove](m, &FilesMove{}); err != nil {
		return fae.Wrap(err, "registering worker: files_move (FilesMove)")
	}

	if err := minion.Register[*FilesRemoveOld](m, &FilesRemoveOld{}); err != nil {
		return fae.Wrap(err, "registering worker: files_remove_old (FilesRemoveOld)")
	}

	if err := minion.Register[*FilesRename](m, &FilesRename{}); err != nil {
		return fae.Wrap(err, "registering worker: files_rename (FilesRename)")
	}

	if err := minion.Register[*FilesRenameMedium](m, &FilesRenameMedium{}); err != nil {
		return fae.Wrap(err, "registering worker: files_rename_medium (FilesRenameMedium)")
	}

	if err := minion.Register[*LibraryCounts](m, &LibraryCounts{}); err != nil {
		return fae.Wrap(err, "registering worker: library_counts (LibraryCounts)")
	}
	if _, err := m.Schedule("0 0 * * * *", &LibraryCounts{}); err != nil {
		return fae.Wrap(err, "scheduling worker: library_counts (LibraryCounts)")
	}

	if err := minion.Register[*MediaImages](m, &MediaImages{}); err != nil {
		return fae.Wrap(err, "registering worker: media_images (MediaImages)")
	}

	if err := minion.Register[*MediumImage](m, &MediumImage{}); err != nil {
		return fae.Wrap(err, "registering worker: medium_image (MediumImage)")
	}

	if err := minion.Register[*MigratePaths](m, &MigratePaths{}); err != nil {
		return fae.Wrap(err, "registering worker: migrate_paths (MigratePaths)")
	}

	if err := minion.Register[*MovieDelete](m, &MovieDelete{}); err != nil {
		return fae.Wrap(err, "registering worker: movie_delete (MovieDelete)")
	}

	if err := minion.Register[*MovieUpdate](m, &MovieUpdate{}); err != nil {
		return fae.Wrap(err, "registering worker: movie_update (MovieUpdate)")
	}

	if err := minion.Register[*MovieUpdateAll](m, &MovieUpdateAll{}); err != nil {
		return fae.Wrap(err, "registering worker: movie_update_all (MovieUpdateAll)")
	}
	if _, err := m.Schedule("0 0 10 * * 0", &MovieUpdateAll{}); err != nil {
		return fae.Wrap(err, "scheduling worker: movie_update_all (MovieUpdateAll)")
	}

	if err := minion.Register[*NzbgetProcess](m, &NzbgetProcess{}); err != nil {
		return fae.Wrap(err, "registering worker: nzbget_process (NzbgetProcess)")
	}

	if err := minion.Register[*PathDelete](m, &PathDelete{}); err != nil {
		return fae.Wrap(err, "registering worker: path_delete (PathDelete)")
	}

	if err := minion.Register[*PathDeleteAll](m, &PathDeleteAll{}); err != nil {
		return fae.Wrap(err, "registering worker: path_delete_all (PathDeleteAll)")
	}

	if err := minion.RegisterWithQueue[*PathImport](m, &PathImport{}, "paths"); err != nil {
		return fae.Wrap(err, "registering worker: path_import (PathImport)")
	}

	if err := minion.Register[*PathManage](m, &PathManage{}); err != nil {
		return fae.Wrap(err, "registering worker: path_manage (PathManage)")
	}

	if err := minion.RegisterWithQueue[*PathManageAll](m, &PathManageAll{}, "paths"); err != nil {
		return fae.Wrap(err, "registering worker: path_manage_all (PathManageAll)")
	}

	if err := minion.Register[*PlexCollectionUpdate](m, &PlexCollectionUpdate{}); err != nil {
		return fae.Wrap(err, "registering worker: plex_collection_update (PlexCollectionUpdate)")
	}

	if err := minion.Register[*PlexFiles](m, &PlexFiles{}); err != nil {
		return fae.Wrap(err, "registering worker: plex_files (PlexFiles)")
	}
	if _, err := m.Schedule("0 0 * * * *", &PlexFiles{}); err != nil {
		return fae.Wrap(err, "scheduling worker: plex_files (PlexFiles)")
	}

	if err := minion.Register[*PlexFilesPartial](m, &PlexFilesPartial{}); err != nil {
		return fae.Wrap(err, "registering worker: plex_files_partial (PlexFilesPartial)")
	}

	if err := minion.Register[*PlexLibraryShoworder](m, &PlexLibraryShoworder{}); err != nil {
		return fae.Wrap(err, "registering worker: plex_library_showorder (PlexLibraryShoworder)")
	}
	if _, err := m.Schedule("0 0 10 * * *", &PlexLibraryShoworder{}); err != nil {
		return fae.Wrap(err, "scheduling worker: plex_library_showorder (PlexLibraryShoworder)")
	}

	if err := minion.Register[*PlexPinToUsers](m, &PlexPinToUsers{}); err != nil {
		return fae.Wrap(err, "registering worker: plex_pin_to_users (PlexPinToUsers)")
	}

	if err := minion.Register[*PlexUserUpdates](m, &PlexUserUpdates{}); err != nil {
		return fae.Wrap(err, "registering worker: plex_user_updates (PlexUserUpdates)")
	}

	if err := minion.Register[*PlexWatched](m, &PlexWatched{}); err != nil {
		return fae.Wrap(err, "registering worker: plex_watched (PlexWatched)")
	}
	if _, err := m.Schedule("0 0 * * * *", &PlexWatched{}); err != nil {
		return fae.Wrap(err, "scheduling worker: plex_watched (PlexWatched)")
	}

	if err := minion.Register[*PlexWatchedAll](m, &PlexWatchedAll{}); err != nil {
		return fae.Wrap(err, "registering worker: plex_watched_all (PlexWatchedAll)")
	}
	if _, err := m.Schedule("0 0 11 * * 0", &PlexWatchedAll{}); err != nil {
		return fae.Wrap(err, "scheduling worker: plex_watched_all (PlexWatchedAll)")
	}

	if err := minion.Register[*PlexWatchlistUpdates](m, &PlexWatchlistUpdates{}); err != nil {
		return fae.Wrap(err, "registering worker: plex_watchlist_updates (PlexWatchlistUpdates)")
	}
	if _, err := m.Schedule("0 0 * * * *", &PlexWatchlistUpdates{}); err != nil {
		return fae.Wrap(err, "scheduling worker: plex_watchlist_updates (PlexWatchlistUpdates)")
	}

	if err := minion.Register[*PlexWebhook](m, &PlexWebhook{}); err != nil {
		return fae.Wrap(err, "registering worker: plex_webhook (PlexWebhook)")
	}

	if err := minion.Register[*ResetIndexes](m, &ResetIndexes{}); err != nil {
		return fae.Wrap(err, "registering worker: reset_indexes (ResetIndexes)")
	}

	if err := minion.Register[*SeriesDelete](m, &SeriesDelete{}); err != nil {
		return fae.Wrap(err, "registering worker: series_delete (SeriesDelete)")
	}

	if err := minion.RegisterWithQueue[*SeriesUpdate](m, &SeriesUpdate{}, "series"); err != nil {
		return fae.Wrap(err, "registering worker: series_update (SeriesUpdate)")
	}

	if err := minion.Register[*SeriesUpdateAll](m, &SeriesUpdateAll{}); err != nil {
		return fae.Wrap(err, "registering worker: series_update_all (SeriesUpdateAll)")
	}
	if _, err := m.Schedule("0 0 10 * * 0", &SeriesUpdateAll{}); err != nil {
		return fae.Wrap(err, "scheduling worker: series_update_all (SeriesUpdateAll)")
	}

	if err := minion.Register[*SeriesUpdateDonghua](m, &SeriesUpdateDonghua{}); err != nil {
		return fae.Wrap(err, "registering worker: series_update_donghua (SeriesUpdateDonghua)")
	}
	if _, err := m.Schedule("0 0 0 * * *", &SeriesUpdateDonghua{}); err != nil {
		return fae.Wrap(err, "scheduling worker: series_update_donghua (SeriesUpdateDonghua)")
	}

	if err := minion.Register[*SeriesUpdateKind](m, &SeriesUpdateKind{}); err != nil {
		return fae.Wrap(err, "registering worker: series_update_kind (SeriesUpdateKind)")
	}

	if err := minion.Register[*SeriesUpdateRecent](m, &SeriesUpdateRecent{}); err != nil {
		return fae.Wrap(err, "registering worker: series_update_recent (SeriesUpdateRecent)")
	}
	if _, err := m.Schedule("0 */15 * * * *", &SeriesUpdateRecent{}); err != nil {
		return fae.Wrap(err, "scheduling worker: series_update_recent (SeriesUpdateRecent)")
	}

	if err := minion.Register[*SeriesUpdateToday](m, &SeriesUpdateToday{}); err != nil {
		return fae.Wrap(err, "registering worker: series_update_today (SeriesUpdateToday)")
	}
	if _, err := m.Schedule("0 0 */12 * * *", &SeriesUpdateToday{}); err != nil {
		return fae.Wrap(err, "scheduling worker: series_update_today (SeriesUpdateToday)")
	}

	if err := minion.Register[*UpdateIndexes](m, &UpdateIndexes{}); err != nil {
		return fae.Wrap(err, "registering worker: update_indexes (UpdateIndexes)")
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
