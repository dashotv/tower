package app

import (
	"github.com/dashotv/minion"
)

// func init() {
// 	starters = append(starters, func(_ context.Context, a *Application) error {
// 		a.Workers.Subscribe(a.MinionNotification)
// 		return nil
// 	})
// }

var workersList = map[string]minion.Payload{
	"DownloadsProcess": &DownloadsProcess{},

	"CleanupLogs":             &CleanupLogs{},
	"CleanPlexPins":           &CleanPlexPins{},
	"PlexPinToUsers":          &PlexPinToUsers{},
	"PlexUserUpdates":         &PlexUserUpdates{},
	"PlexWatchlistUpdates":    &PlexWatchlistUpdates{},
	"CreateMediaFromRequests": &CreateMediaFromRequests{},

	"UpdateIndexes": &UpdateIndexes{},
	// "DownloadsFileMove":        &DownloadFileMover{},

	"FileWalk":  &FileWalk{},
	"FileMatch": &FileMatch{},

	"PathCleanupAll": &PathCleanupAll{},

	"SeriesUpdateAll":     &SeriesUpdateAll{},
	"SeriesUpdateDonghua": &SeriesUpdateKind{SeriesKind: "donghua"},
}

// func (a *Application) MinionNotification(n *minion.Notification) {
// 	a.Log.Named("minion.notification").Infof("Received notification: %+v", n)
// 	if n.Event == "job:success" || n.Event == "job:failure" {
// 		a.Events.Send("tower.jobs", &EventJobs{Event: n.Event, ID: n.JobID, Kind: n.Kind})
// 	}
// }
