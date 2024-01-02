package app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/minion"
	"github.com/dashotv/tvdb"
)

// TvdbUpdateSeriesEpisodes
type TvdbUpdateSeriesEpisodes struct {
	minion.WorkerDefaults[*TvdbUpdateSeriesEpisodes]
	ID string
}

func (j *TvdbUpdateSeriesEpisodes) Kind() string { return "TvdbUpdateSeriesEpisodes" }
func (j *TvdbUpdateSeriesEpisodes) Work(ctx context.Context, job *minion.Job[*TvdbUpdateSeriesEpisodes]) error {
	// log :=app.Log.Named("TvdbUpdateSeriesEpisodes")
	//app.Log.Info("updating series episodes")

	id := job.Args.ID

	series := &Series{}
	err := app.DB.Series.Find(id, series)
	if err != nil {
		return errors.Wrap(err, "finding series")
	}

	sid, err := strconv.Atoi(series.SourceId)
	if err != nil {
		return errors.Wrap(err, "converting source id")
	}

	// resp, err := app.Tvdb.GetSeriesExtended(int64(sid), operations.GetSeriesExtendedMetaEpisodes.ToPointer(), tvdb.Bool(true))
	resp, err := app.Tvdb.GetSeriesSeasonEpisodesTranslated(int64(sid), "eng", 0, "default")
	if err != nil {
		return errors.Wrap(err, "getting episodes")
	}

	if resp.Data == nil {
		return errors.New("no data")
	}

	episodeMap, err := buildEpisodeMap(id)
	if err != nil {
		return errors.Wrap(err, "building episode map")
	}
	//app.Log.Infof("episode map: %d", len(episodeMap))
	//app.Log.Infof("episodes: %d", len(resp.Data.Episodes))

	for _, e := range resp.Data.Episodes {
		episode := episodeMap[tvdb.Int64Value(e.SeasonNumber)][tvdb.Int64Value(e.Number)]
		if episode == nil {
			episode = &Episode{}
		}

		//app.Log.Infof("creating/updating episode %d/%d %s", tvdb.Int64Value(e.SeasonNumber), tvdb.Int64Value(e.Number), tvdb.StringValue(e.Aired))
		episode.Type = "Episode"
		episode.SeriesId = series.ID
		episode.SourceId = fmt.Sprintf("%d", tvdb.Int64Value(e.ID))
		// episode.AbsoluteNumber = int(tvdb.Int64Value(e.AbsoluteNumber))
		episode.SeasonNumber = int(tvdb.Int64Value(e.SeasonNumber))
		episode.EpisodeNumber = int(tvdb.Int64Value(e.Number))
		episode.Title = tvdb.StringValue(e.Name)
		episode.Description = tvdb.StringValue(e.Overview)
		if tvdb.StringValue(e.Aired) != "" {
			date, err := time.Parse("2006-01-02", tvdb.StringValue(e.Aired))
			if err != nil {
				return errors.Wrap(err, "parsing release date")
			}
			episode.ReleaseDate = date
		} else {
			episode.ReleaseDate = time.Unix(0, 0)
		}

		if series.Kind == "anime" {
			resp, err := app.Tvdb.GetEpisodeExtended(tvdb.Int64Value(e.ID), nil)
			if err != nil {
				return errors.Wrap(err, "getting episodes")
			}
			if resp.Data == nil {
				return errors.New("no episode data")
			}
			for _, ee := range resp.Data.Seasons {
				if tvdb.StringValue(ee.Type.Type) == "absolute" {
					episode.AbsoluteNumber = int(tvdb.Int64Value(resp.Data.Number))
					break
				}
			}
		}

		if err := app.DB.Episode.Save(episode); err != nil {
			return errors.Wrap(err, fmt.Sprintf("updating episode %s %d/%d", id, episode.SeasonNumber, episode.EpisodeNumber))
		}
	}

	return nil
}

func buildEpisodeMap(id string) (map[int64]map[int64]*Episode, error) {
	episodeMap := map[int64]map[int64]*Episode{}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(err, "converting id")
	}

	episodes, err := app.DB.Episode.Query().Where("series_id", oid).Limit(-1).Run()
	if err != nil {
		return nil, errors.Wrap(err, "querying episodes")
	}

	app.Log.Warnf("episodes: %d", len(episodes))

	for _, e := range episodes {
		sn := int64(e.SeasonNumber)
		en := int64(e.EpisodeNumber)
		if episodeMap[sn] == nil {
			episodeMap[sn] = map[int64]*Episode{}
		}
		episodeMap[sn][en] = e
	}

	return episodeMap, nil
}
