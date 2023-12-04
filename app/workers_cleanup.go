package app

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/minion"
)

type CleanupLogs struct{}

func (j *CleanupLogs) Kind() string { return "CleanupLogs" }
func (j *CleanupLogs) Work(ctx context.Context, job *minion.Job[*CleanupLogs]) error {
	_, err := db.Message.Collection.DeleteMany(context.Background(), bson.M{"created_at": bson.M{"$lt": time.Now().UTC().AddDate(0, 0, -5)}})
	if err != nil {
		return errors.Wrap(err, "cleaning logs")
	}
	return nil
}

type CleanupJobs struct{}

func (j *CleanupJobs) Kind() string { return "CleanupJobs" }
func (j *CleanupJobs) Work(ctx context.Context, job *minion.Job[*CleanupJobs]) error {
	_, err := db.Minion.Collection.DeleteMany(context.Background(), bson.M{"created_at": bson.M{"$lt": time.Now().UTC().AddDate(0, 0, -1)}})
	if err != nil {
		return errors.Wrap(err, "cleaning logs")
	}
	_, err = db.Minion.Collection.DeleteMany(context.Background(), bson.M{"kind": "ping"})
	if err != nil {
		return errors.Wrap(err, "cleaning logs")
	}

	return nil
}
