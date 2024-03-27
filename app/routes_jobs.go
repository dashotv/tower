package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/minion/database"
)

func (a *Application) JobsIndex(c echo.Context, page int, limit int) error {
	if limit == 0 {
		limit = 25
	}
	skip := 0
	if page > 0 {
		skip = (page - 1) * limit
	}

	status := c.QueryParam("status")
	q := app.DB.Minion.Query()
	if status != "" {
		q = q.Where("status", status)
	}

	count, err := q.Count()
	if err != nil {
		return err
	}

	list, err := q.Skip(skip).Limit(limit).Desc("updated_at").Run()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, H{"total": count, "jobs": list})
}

func (a *Application) JobsCreate(c echo.Context, job string) error {
	a.Log.Debugf("JobsCreate: %s", job)
	if job == "" {
		return errors.New("missing job")
	}

	j, ok := workersList[job]
	if !ok || j == nil {
		return fmt.Errorf("unknown job: %s", job)
	}

	app.Log.Infof("Enqueuing job: %s", j.Kind())
	err := app.Workers.Enqueue(j)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, j)
}

func (a *Application) JobsDelete(c echo.Context, id string, hard bool) error {
	if id == string(database.StatusPending) && !hard {
		filter := bson.M{"status": database.StatusPending}
		if _, err := app.DB.Minion.Collection.UpdateMany(context.Background(), filter, bson.M{"$set": bson.M{"status": database.StatusCancelled}}); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, H{"error": false})
	} else if id == string(database.StatusFailed) && hard {
		filter := bson.M{"status": database.StatusFailed}
		if _, err := app.DB.Minion.Collection.DeleteMany(context.Background(), filter); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, H{"error": false})
	} else if id == string(database.StatusCancelled) && hard {
		filter := bson.M{"status": database.StatusCancelled}
		if _, err := app.DB.Minion.Collection.DeleteMany(context.Background(), filter); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, H{"error": false})
	}

	j, err := app.DB.Minion.Get(id, &Minion{})
	if err != nil {
		return err
	}

	j.Status = string(database.StatusCancelled)
	if err := app.DB.Minion.Save(j); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, H{"error": false})
}
