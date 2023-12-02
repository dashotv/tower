package app

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/minion"
)

type CleanupLogs struct{}

func (j *CleanupLogs) Kind() string { return "cleanup_logs" }
func (j *CleanupLogs) Work(ctx context.Context, job *minion.Job[*CleanupLogs]) error {
	_, err := db.Message.Collection.DeleteMany(context.Background(), bson.M{"created_at": bson.M{"$lt": time.Now().UTC().AddDate(0, 0, -5)}})
	if err != nil {
		return errors.Wrap(err, "cleaning logs")
	}

	return nil
}

type CleanupJobs struct{}

func (j *CleanupJobs) Kind() string { return "cleanup_jobs" }
func (j *CleanupJobs) Work(ctx context.Context, job *minion.Job[*CleanupJobs]) error {
	list, err := db.Minion.Query().
		GreaterThan("created_at", time.Now().UTC().AddDate(0, 0, -1)).
		Run()
	if err != nil {
		return errors.Wrap(err, "querying jobs")
	}

	for _, j := range list {
		err := db.Minion.Delete(j)
		if err != nil {
			return errors.Wrap(err, "deleting job")
		}
	}

	return nil
}
