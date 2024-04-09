package app

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
	"github.com/dashotv/tower/internal/importer"
)

type SeriesUpdateAll struct {
	minion.WorkerDefaults[*SeriesUpdateAll]
}

func (j *SeriesUpdateAll) Kind() string { return "series_update_all" }
func (j *SeriesUpdateAll) Work(ctx context.Context, job *minion.Job[*SeriesUpdateAll]) error {
	q := app.DB.Series.Query().LessThan("updated_at", time.Now().Add(-24*time.Hour*7))
	total, err := q.Count()
	if err != nil {
		return err
	}

	for skip := 0; skip < int(total); skip += 100 {
		list, err := q.Skip(skip).Limit(100).Run()
		if err != nil {
			return err
		}

		for _, series := range list {
			if err := app.Workers.Enqueue(&SeriesUpdate{ID: series.ID.Hex(), SkipImages: true, Title: series.Title}); err != nil {
				return err
			}
		}
	}

	return nil
}

type SeriesUpdateKind struct {
	minion.WorkerDefaults[*SeriesUpdateKind]
	SeriesKind string
}

func (j *SeriesUpdateKind) Kind() string { return "SeriesUpdateKind" }
func (j *SeriesUpdateKind) Work(ctx context.Context, job *minion.Job[*SeriesUpdateKind]) error {
	q := app.DB.Series.Query().Where("kind", job.Args.SeriesKind)
	total, err := q.Count()
	if err != nil {
		return err
	}

	for skip := 0; skip < int(total); skip += 100 {
		list, err := q.Skip(skip).Limit(100).Run()
		if err != nil {
			return err
		}

		for _, series := range list {
			if err := app.Workers.Enqueue(&SeriesUpdate{ID: series.ID.Hex(), SkipImages: true, Title: series.Title}); err != nil {
				return err
			}
		}
	}

	return nil
}

type SeriesUpdateDonghua struct {
	minion.WorkerDefaults[*SeriesUpdateDonghua]
}

func (j *SeriesUpdateDonghua) Kind() string { return "series_update_donghua" }
func (j *SeriesUpdateDonghua) Work(ctx context.Context, job *minion.Job[*SeriesUpdateDonghua]) error {
	//args := job.Args
	return app.Workers.Enqueue(&SeriesUpdateKind{SeriesKind: "donghua"})
}

type SeriesUpdateRecent struct {
	minion.WorkerDefaults[*SeriesUpdateRecent]
}

func (j *SeriesUpdateRecent) Kind() string { return "series_update_recent" }
func (j *SeriesUpdateRecent) Work(ctx context.Context, job *minion.Job[*SeriesUpdateRecent]) error {
	ints, err := app.Importer.SeriesUpdated(time.Now().Add(-15 * time.Minute).Unix())
	if err != nil {
		return fae.Wrap(err, "recent")
	}

	for _, id := range ints {
		list, err := app.DB.Series.Query().Where("source", "tvdb").Where("source_id", fmt.Sprintf("%d", id)).Run()
		if err != nil {
			return fae.Wrap(err, "recent: list")
		}
		for _, series := range list {
			if err := app.Workers.Enqueue(&SeriesUpdate{ID: series.ID.Hex(), SkipImages: true, Title: series.Title}); err != nil {
				return fae.Wrap(err, "recent: enqueue")
			}
		}
	}

	return nil
}

type SeriesUpdateToday struct {
	minion.WorkerDefaults[*SeriesUpdateToday]
}

func (j *SeriesUpdateToday) Kind() string { return "series_update_today" }
func (j *SeriesUpdateToday) Work(ctx context.Context, job *minion.Job[*SeriesUpdateToday]) error {
	today := time.Now().Format("2006-01-02")
	list, err := app.DB.Episode.Query().Where("release_date", today).Run()
	if err != nil {
		return fae.Wrap(err, "listing todays episodes")
	}

	seriesIds := lo.Map(list, func(e *Episode, i int) string {
		return e.SeriesId.Hex()
	})
	seriesIds = lo.Uniq(seriesIds)

	for _, id := range seriesIds {
		if err := app.Workers.Enqueue(&SeriesUpdate{ID: id, SkipImages: true}); err != nil {
			return fae.Wrap(err, "enqueueing series")
		}
	}

	return nil
}

type SeriesUpdate struct {
	minion.WorkerDefaults[*SeriesUpdate]
	ID         string `bson:"id" json:"id"`
	Title      string `bson:"title" json:"title"`
	SkipImages bool   `bson:"skip_images" json:"skip_images"`
}

func (j *SeriesUpdate) Kind() string { return "series_update" }
func (j *SeriesUpdate) Work(ctx context.Context, job *minion.Job[*SeriesUpdate]) error {
	id := job.Args.ID
	eg, ctx := errgroup.WithContext(ctx)

	series := &Series{}
	err := app.DB.Series.Find(id, series)
	if err != nil {
		return err
	}

	if series.Source != "tvdb" {
		return nil
	}

	tvdbid, err := strconv.ParseInt(series.SourceId, 10, 64)
	if err != nil {
		return fae.Wrap(err, "converting source id")
	}

	eg.Go(func() error {
		s, err := app.Importer.Series(tvdbid)
		if err != nil {
			return fae.Wrap(err, "importer series")
		}

		series.Title = s.Title
		series.Description = s.Description
		series.Status = s.Status
		series.ReleaseDate = dateFromString(s.Airdate)
		if series.Display == "" {
			series.Display = s.Title
		}
		if series.Search == "" {
			series.Search = path(s.Title)
		}
		if series.Directory == "" {
			series.Directory = path(s.Title)
		}

		return nil
	})

	eg.Go(func() error {
		order := importer.EpisodeOrderDefault
		anime := isAnimeKind(string(series.Kind))
		if anime {
			order = importer.EpisodeOrderAbsolute
		}

		eps, err := app.Importer.SeriesEpisodes(tvdbid, order)
		if err != nil {
			return fae.Wrap(err, "importer series episodes")
		}

		episodeMap, err := episodeMap(id)
		if err != nil {
			return fae.Wrap(err, "building episode map")
		}

		found := []int64{}

		for _, e := range eps {
			episode, ok := episodeMap[e.ID]
			if ok {
				found = append(found, e.ID)
			}
			if episode == nil {
				episode = &Episode{}
			}

			episode.Type = "Episode"
			episode.SeriesId = series.ID
			episode.SourceId = fmt.Sprintf("%d", e.ID)
			episode.SeasonNumber = e.Season
			episode.EpisodeNumber = e.Episode
			episode.AbsoluteNumber = e.Absolute
			episode.Title = e.Title
			episode.Description = e.Description
			episode.ReleaseDate = dateFromString(e.Airdate)

			if err := app.DB.Episode.Save(episode); err != nil {
				return fae.Wrap(err, fmt.Sprintf("updating episode %s %d/%d", id, episode.SeasonNumber, episode.EpisodeNumber))
			}
		}

		all := lo.Keys(episodeMap)
		missing, updated := lo.Difference(all, found)
		if _, err := app.DB.Episode.Collection.UpdateMany(ctx, bson.M{"_type": "Episode", "series_id": series.ID, "source_id": bson.M{"$in": missing}}, bson.M{"$set": bson.M{"missing": time.Now()}}); err != nil {
			return fae.Wrap(err, "missing")
		}
		if _, err := app.DB.Episode.Collection.UpdateMany(ctx, bson.M{"_type": "Episode", "series_id": series.ID, "source_id": bson.M{"$in": updated}}, bson.M{"$set": bson.M{"missing": nil}}); err != nil {
			return fae.Wrap(err, "found")
		}
		if _, err := app.DB.Episode.Collection.DeleteMany(ctx, bson.M{"_type": "Episode", "series_id": series.ID, "missing": bson.M{"$ne": nil}, "paths.type": bson.M{"$ne": "video"}}); err != nil {
			return fae.Wrap(err, "missing delete")
		}

		// db.media.aggregate([{ $match: { _type: "Episode", series_id: ObjectId("65b572ff28653636fbae17de") } }, { $group: { _id: { s: "$season_number", e: "$episode_number", a: "$absolute_number" }, dups: { $push: '$_id' } } }, { $sort: { dups: -1 } }])
		cur, err := app.DB.Episode.Collection.Aggregate(ctx, bson.A{
			bson.M{"$match": bson.M{"_type": "Episode", "series_id": series.ID}},
			bson.M{"$group": bson.M{
				"_id":  bson.M{"s": "$season_number", "e": "$episode_number", "a": "$absolute_number"},
				"dups": bson.M{"$push": "$_id"}},
			},
			bson.M{"$sort": bson.M{"dups": -1}}})
		if err != nil {
			return fae.Wrap(err, "duplicates")
		}
		for cur.Next(ctx) {
			result := &SeriesDupResult{}
			if err := cur.Decode(result); err != nil {
				return fae.Wrap(err, "decoding")
			}
			if len(result.Dups) > 1 {
				for _, id := range result.Dups[1:] {
					if _, err := app.DB.Episode.Collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
						return fae.Wrap(err, "deleting")
					}
				}
			}
		}

		return nil
	})

	eg.Go(func() error {
		if !job.Args.SkipImages {
			covers, backgrounds, err := app.Importer.SeriesImages(tvdbid)
			if err != nil {
				return fae.Wrap(err, "importer series images")
			}

			if len(covers) > 0 {
				eg.Go(func() error {
					err := seriesImage(series, "cover", covers[0], posterRatio)
					if err != nil {
						app.Log.Errorf("series %s cover: %v", series.ID.Hex(), err)
					}
					return nil
				})
			}
			if len(backgrounds) > 0 {
				eg.Go(func() error {
					seriesImage(series, "background", backgrounds[0], backgroundRatio)
					if err != nil {
						app.Log.Errorf("series %s background: %v", series.ID.Hex(), err)
					}
					return nil
				})
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return fae.Wrapf(err, "series: %s", series.Title)
	}

	if err := app.DB.Series.Save(series); err != nil {
		return fae.Wrapf(err, "saving series: %s", series.Title)
	}
	return nil
}

type SeriesImage struct {
	minion.WorkerDefaults[*SeriesImage]
	ID    string
	Type  string
	Path  string
	Ratio float32
}

func (j *SeriesImage) Kind() string { return "SeriesImage" }
func (j *SeriesImage) Work(ctx context.Context, job *minion.Job[*SeriesImage]) error {
	id := job.Args.ID
	t := job.Args.Type
	remote := job.Args.Path
	ratio := job.Args.Ratio

	series := &Series{}
	if err := app.DB.Series.Find(id, series); err != nil {
		return fae.Wrap(err, "finding series")
	}

	return seriesImage(series, t, remote, ratio)
}

// TODO: make this a function
func seriesImage(series *Series, t string, remote string, ratio float32) error {
	extension := filepath.Ext(remote)
	if len(extension) > 0 && extension[0] == '.' {
		extension = extension[1:]
	}
	local := fmt.Sprintf("series-%s/%s", series.ID.Hex(), t)
	dest := fmt.Sprintf("%s/%s.%s", app.Config.DirectoriesImages, local, extension)
	thumb := fmt.Sprintf("%s/%s_thumb.%s", app.Config.DirectoriesImages, local, extension)

	if err := imageDownload(remote, dest); err != nil {
		return fae.Wrap(err, "downloading image")
	}

	height := 400
	width := int(float32(height) * ratio)
	if err := imageResize(dest, thumb, width, height); err != nil {
		return fae.Wrap(err, "resizing image")
	}

	var img *Path
	for _, p := range series.Paths {
		if string(p.Type) == t {
			img = p
			break
		}
	}

	if img == nil {
		img = &Path{}
		series.Paths = append(series.Paths, img)
	}

	img.Type = primitive.Symbol(t)
	img.Remote = remote
	img.Local = local
	img.Extension = extension

	if err := app.DB.Series.Update(series); err != nil {
		return fae.Wrap(err, "updating series")
	}

	return nil
}

func dateFromString(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Unix(0, 0)
	}
	return t
}

func episodeMap(id string) (map[int64]*Episode, error) {
	episodeMap := make(map[int64]*Episode)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fae.Wrap(err, "converting id")
	}

	episodes, err := app.DB.Episode.Query().Where("series_id", oid).Limit(-1).Run()
	if err != nil {
		return nil, fae.Wrap(err, "querying episodes")
	}

	for _, e := range episodes {
		sid, err := strconv.ParseInt(e.SourceId, 10, 64)
		if err != nil {
			return nil, fae.Wrap(err, "converting source id")
		}
		episodeMap[sid] = e
	}

	return episodeMap, nil
}

func episodeSEMap(id string, anime bool) (map[int]map[int]*Episode, error) {
	episodeMap := map[int]map[int]*Episode{}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fae.Wrap(err, "converting id")
	}

	episodes, err := app.DB.Episode.Query().Where("series_id", oid).Limit(-1).Run()
	if err != nil {
		return nil, fae.Wrap(err, "querying episodes")
	}

	app.Log.Warnf("episodes: %d", len(episodes))

	for _, e := range episodes {
		sn := e.SeasonNumber
		en := e.EpisodeNumber
		if anime {
			sn = 1
			en = e.AbsoluteNumber
		}
		if episodeMap[sn] == nil {
			episodeMap[sn] = map[int]*Episode{}
		}
		episodeMap[sn][en] = e
	}

	return episodeMap, nil
}

type SeriesDupResult struct {
	Season   int                  `bson:"_id.s" json:"season"`
	Episode  int                  `bson:"_id.e" json:"episode"`
	Absolute int                  `bson:"_id.a" json:"absolute"`
	Dups     []primitive.ObjectID `bson:"dups" json:"dups"`
}