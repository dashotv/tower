package app

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
)

// CleanPlexPins removes old pins
type CleanPlexPins struct {
	minion.WorkerDefaults[*CleanPlexPins]
}

func (j *CleanPlexPins) Kind() string { return "CleanPlexPins" }
func (j *CleanPlexPins) Work(ctx context.Context, job *minion.Job[*CleanPlexPins]) error {
	list, err := app.DB.Pin.Query().
		GreaterThan("created_at", time.Now().UTC().AddDate(0, 0, -1)).
		Run()
	if err != nil {
		return fae.Wrap(err, "querying pins")
	}

	for _, p := range list {
		err := app.DB.Pin.Delete(p)
		if err != nil {
			return fae.Wrap(err, "deleting pin")
		}
	}

	return nil
}

type CleanupLogs struct {
	minion.WorkerDefaults[*CleanupLogs]
}

func (j *CleanupLogs) Kind() string { return "cleanup_logs" }
func (j *CleanupLogs) Work(ctx context.Context, job *minion.Job[*CleanupLogs]) error {
	if _, err := app.DB.Message.Collection.DeleteMany(context.Background(), bson.M{"created_at": bson.M{"$lt": time.Now().UTC().AddDate(0, 0, -3)}}); err != nil {
		return fae.Wrap(err, "deleting messages")
	}
	return nil
}
