package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/dashotv/minion"
)

func (a *Application) JobsIndex(c echo.Context, page int, limit int) error {
	if limit == 0 {
		limit = 25
	}
	skip := (page * limit) - limit
	list, err := app.DB.Minion.Query().Skip(skip).Limit(limit).Desc("created_at").Run()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"jobs": list})
}

func (a *Application) JobsCreate(c echo.Context, job string) error {
	j := workersList[job]
	if j == nil {
		return errors.New("invalid job: " + job)
	}

	app.Log.Infof("Enqueuing job: %s", j.Kind())
	err := app.Workers.Enqueue(j)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, j)
}

func (a *Application) JobsDelete(c echo.Context, id string, hard bool) error {
	if id == string(minion.StatusPending) && !hard {
		filter := bson.M{"status": minion.StatusPending}
		if _, err := app.DB.Minion.Collection.UpdateMany(context.Background(), filter, bson.M{"$set": bson.M{"status": minion.StatusCancelled}}); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, gin.H{"error": false})
	} else if id == string(minion.StatusFailed) && hard {
		filter := bson.M{"status": minion.StatusFailed}
		if _, err := app.DB.Minion.Collection.DeleteMany(context.Background(), filter); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, gin.H{"error": false})
	} else if id == string(minion.StatusCancelled) && hard {
		filter := bson.M{"status": minion.StatusCancelled}
		if _, err := app.DB.Minion.Collection.DeleteMany(context.Background(), filter); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, gin.H{"error": false})
	}

	j, err := app.DB.Minion.Get(id, &Minion{})
	if err != nil {
		return err
	}

	j.Status = string(minion.StatusCancelled)
	if err := app.DB.Minion.Save(j); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{"error": false})
}
