package app

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/dashotv/minion"
	"github.com/dashotv/tvdb"
)

// TvdbUpdateSeries
type TvdbUpdateSeries struct {
	minion.WorkerDefaults[*TvdbUpdateSeries]
	ID        string
	JustMedia bool
}

func (j *TvdbUpdateSeries) Kind() string { return "TvdbUpdateSeries" }
func (j *TvdbUpdateSeries) Work(ctx context.Context, job *minion.Job[*TvdbUpdateSeries]) error {
	id := job.Args.ID

	series := &Series{}
	err := app.DB.Series.Find(id, series)
	if err != nil {
		return err
	}

	sid, err := strconv.Atoi(series.SourceId)
	if err != nil {
		return errors.Wrap(err, "converting source id")
	}

	{
		resp, err := app.Tvdb.GetSeriesTranslation(int64(sid), "eng")
		if err != nil {
			return err
		}

		if resp.Data == nil {
			return errors.New("no data")
		}

		series.Title = tvdb.StringValue(resp.Data.Name)
		if series.Display == "" {
			series.Display = series.Title
		}
		if series.Search == "" {
			series.Search = path(series.Title)
		}
		if series.Directory == "" {
			series.Directory = path(series.Title)
		}
		series.Description = tvdb.StringValue(resp.Data.Overview)
	}

	resp, err := app.Tvdb.GetSeriesBase(int64(sid))
	if err != nil {
		return err
	}

	if resp.Data == nil {
		return errors.New("no data")
	}

	data := resp.Data
	series.Status = tvdb.StringValue(data.Status.Name)

	date, err := time.Parse("2006-01-02", tvdb.StringValue(data.FirstAired))
	if err != nil {
		return errors.Wrap(err, "parsing release date")
	}
	series.ReleaseDate = date

	if err := app.DB.Series.Update(series); err != nil {
		return errors.Wrap(err, "updating series")
	}
	if !job.Args.JustMedia {
		if err := TvdbUpdateSeriesCover(series.ID.Hex(), int64(sid)); err != nil {
			app.Log.Warnf("failed to update cover: %s", err)
		}
		if err := TvdbUpdateSeriesBackground(series.ID.Hex(), int64(sid)); err != nil {
			app.Log.Warnf("failed to update background: %s", err)
		}
		if err := app.Workers.Enqueue(&TvdbUpdateSeriesEpisodes{ID: series.ID.Hex()}); err != nil {
			return errors.Wrap(err, "enqueuing series episodes")
		}
		// if err := app.Workers.Enqueue(&MediaPaths{ID: id}); err != nil {
		// 	return errors.Wrap(err, "enqueuing media paths")
		// }
	}

	return nil
}
