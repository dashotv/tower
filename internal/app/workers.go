package app

import (
	"github.com/dashotv/minion"
)

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
