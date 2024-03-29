package app

import (
	"github.com/dashotv/minion"
)

var workersList = map[string]minion.Payload{
	"DownloadsProcess": &DownloadsProcess{},

	"CleanupLogs":             &CleanupLogs{},
	"CleanupJobs":             &CleanupJobs{},
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

// This allows you to notify other services as jobs change status.
func (a *Application) MinionNotification(n *minion.Notification) {
	if n.JobID == "-" {
		return
	}

	j := &Minion{}
	err := app.DB.Minion.Find(n.JobID, j)
	if err != nil {
		a.Log.Errorf("finding job: %s", err)
		return
	}

	if n.Event == "job:created" {
		a.Events.Send("tower.jobs", &EventJobs{"created", j.ID.Hex(), j})
		return
	}
	a.Events.Send("tower.jobs", &EventJobs{"updated", j.ID.Hex(), j})
}

func (a *Application) MinionStats(stats minion.Stats) {
	a.Events.Send("tower.stats", &stats)
}
