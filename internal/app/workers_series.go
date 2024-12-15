package app

import (
	"context"
	"fmt"
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

type SeriesDelete struct {
	minion.WorkerDefaults[*SeriesDelete]
	ID string `bson:"id" json:"id"`
}

func (j *SeriesDelete) Kind() string { return "series_delete" }
func (j *SeriesDelete) Work(ctx context.Context, job *minion.Job[*SeriesDelete]) error {
	a := ContextApp(ctx)
	id := job.Args.ID

	series := &Series{}
	if err := a.DB.Series.Find(id, series); err != nil {
		return fae.Wrap(err, "finding series")
	}

	// delete files
	if err := a.DB.mediumIdDeletePaths(id); err != nil {
		return fae.Wrap(err, "deleting paths")
	}

	// find episodes
	list, err := a.DB.Episode.Query().Where("series_id", series.ID).Run()
	if err != nil {
		return fae.Wrap(err, "listing episodes")
	}

	// get episode ids
	eids := lo.Map(list, func(e *Episode, i int) primitive.ObjectID {
		return e.ID
	})
	eids = append(eids, series.ID)

	// remove downloads referencing episodes or series
	_, err = a.DB.Download.Collection.DeleteMany(ctx, bson.M{"medium_id": bson.M{"$in": eids}})
	if err != nil {
		return fae.Wrap(err, "deleting downloads")
	}

	// remove watches referencing episodes or series
	_, err = a.DB.Watch.Collection.DeleteMany(ctx, bson.M{"medium_id": bson.M{"$in": eids}})
	if err != nil {
		return fae.Wrap(err, "deleting watches")
	}

	// remove episodes
	_, err = a.DB.Episode.Collection.DeleteMany(ctx, bson.M{"_type": "Episode", "series_id": series.ID})
	if err != nil {
		return fae.Wrap(err, "deleting episodes")
	}

	// remove series
	if err := a.DB.Series.Delete(series); err != nil {
		return fae.Wrap(err, "delete medium")
	}

	return nil
}

type SeriesUpdateAll struct {
	minion.WorkerDefaults[*SeriesUpdateAll]
}

func (j *SeriesUpdateAll) Kind() string { return "series_update_all" }
func (j *SeriesUpdateAll) Work(ctx context.Context, job *minion.Job[*SeriesUpdateAll]) error {
	a := ContextApp(ctx)
	err := a.DB.Series.Query().LessThan("updated_at", time.Now().Add(-24*time.Hour*7)).Batch(100, func(list []*Series) error {
		for _, series := range list {
			if err := a.Workers.Enqueue(&SeriesUpdate{ID: series.ID.Hex(), SkipImages: true, Title: series.Title}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fae.Wrap(err, "batching series")
	}

	return nil
}

type SeriesUpdateKind struct {
	minion.WorkerDefaults[*SeriesUpdateKind]
	SeriesKind string
}

func (j *SeriesUpdateKind) Kind() string { return "SeriesUpdateKind" }
func (j *SeriesUpdateKind) Work(ctx context.Context, job *minion.Job[*SeriesUpdateKind]) error {
	a := ContextApp(ctx)
	q := a.DB.Series.Query().Where("kind", job.Args.SeriesKind)
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
			if err := a.Workers.Enqueue(&SeriesUpdate{ID: series.ID.Hex(), SkipImages: true, Title: series.Title}); err != nil {
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
	a := ContextApp(ctx)
	return a.Workers.Enqueue(&SeriesUpdateKind{SeriesKind: "donghua"})
}

type SeriesUpdateRecent struct {
	minion.WorkerDefaults[*SeriesUpdateRecent]
}

func (j *SeriesUpdateRecent) Kind() string { return "series_update_recent" }
func (j *SeriesUpdateRecent) Work(ctx context.Context, job *minion.Job[*SeriesUpdateRecent]) error {
	a := ContextApp(ctx)
	ints, err := a.Importer.SeriesUpdated(time.Now().Add(-15 * time.Minute).Unix())
	if err != nil {
		return fae.Wrap(err, "recent")
	}

	ints = lo.Uniq(ints)

	for _, id := range ints {
		list, err := a.DB.Series.Query().Where("source", "tvdb").Where("source_id", fmt.Sprintf("%d", id)).Run()
		if err != nil {
			return fae.Wrap(err, "recent: list")
		}
		for _, series := range list {
			if err := a.Workers.Enqueue(&SeriesUpdate{ID: series.ID.Hex(), SkipImages: true, Title: series.Title}); err != nil {
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
	a := ContextApp(ctx)
	today := time.Now().Format("2006-01-02")
	list, err := a.DB.Episode.Query().Where("release_date", today).Run()
	if err != nil {
		return fae.Wrap(err, "listing todays episodes")
	}

	seriesIds := lo.Map(list, func(e *Episode, i int) string {
		return e.SeriesID.Hex()
	})
	seriesIds = lo.Uniq(seriesIds)

	for _, id := range seriesIds {
		if err := a.Workers.Enqueue(&SeriesUpdate{ID: id, SkipImages: true}); err != nil {
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
	a := ContextApp(ctx)
	id := job.Args.ID
	eg, ctx := errgroup.WithContext(ctx)

	series := &Medium{}
	err := a.DB.Medium.Find(id, series)
	if err != nil {
		return err
	}

	if series.Source != "tvdb" || series.SourceID == "" {
		return fae.New("series source not tvdb or missing id")
	}

	tvdbid, err := strconv.ParseInt(series.SourceID, 10, 64)
	if err != nil {
		return fae.Wrap(err, "converting source id")
	}

	eg.Go(func() error {
		s, err := a.Importer.Series(tvdbid)
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

		eps, err := a.Importer.SeriesEpisodes(tvdbid, order)
		if err != nil {
			return fae.Wrap(err, "importer series episodes")
		}

		episodeMap, err := episodeMap(id)
		if err != nil {
			return fae.Wrap(err, "building episode map")
		}

		found := []int64{}

		for _, e := range eps {
			episode := episodeMap.Find(e.ID, e.Season, e.Episode, e.Absolute)
			if episode != nil {
				found = append(found, e.ID)
			} else {
				episode = &Episode{}
			}

			episode.Type = "Episode"
			episode.SeriesID = series.ID
			episode.SourceID = fmt.Sprintf("%d", e.ID)
			episode.SeasonNumber = e.Season
			episode.EpisodeNumber = e.Episode
			episode.AbsoluteNumber = e.Absolute
			episode.Title = e.Title
			episode.Description = e.Description
			episode.ReleaseDate = dateFromString(e.Airdate)

			if err := a.DB.Episode.Save(episode); err != nil {
				return fae.Wrap(err, fmt.Sprintf("updating episode %s %d/%d", id, episode.SeasonNumber, episode.EpisodeNumber))
			}
		}

		all := lo.Keys(episodeMap.byID)
		missing, updated := lo.Difference(all, found)
		if _, err := a.DB.Episode.Collection.UpdateMany(ctx, bson.M{"_type": "Episode", "series_id": series.ID, "source_id": bson.M{"$in": missing}}, bson.M{"$set": bson.M{"missing": time.Now()}}); err != nil {
			return fae.Wrap(err, "missing")
		}
		if _, err := a.DB.Episode.Collection.UpdateMany(ctx, bson.M{"_type": "Episode", "series_id": series.ID, "source_id": bson.M{"$in": updated}}, bson.M{"$set": bson.M{"missing": nil}}); err != nil {
			return fae.Wrap(err, "found")
		}
		if _, err := a.DB.Episode.Collection.DeleteMany(ctx, bson.M{"_type": "Episode", "series_id": series.ID, "missing": bson.M{"$ne": nil}, "paths.type": bson.M{"$ne": "video"}}); err != nil {
			return fae.Wrap(err, "missing delete")
		}

		// db.media.aggregate([{ $match: { _type: "Episode", series_id: ObjectID("65b572ff28653636fbae17de") } }, { $group: { _id: { s: "$season_number", e: "$episode_number", a: "$absolute_number" }, dups: { $push: '$_id' } } }, { $sort: { dups: -1 } }])
		cur, err := a.DB.Episode.Collection.Aggregate(ctx, bson.A{
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
					if _, err := a.DB.Episode.Collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
						return fae.Wrap(err, "deleting")
					}
				}
			}
		}

		return nil
	})

	eg.Go(func() error {
		if !job.Args.SkipImages {
			covers, backgrounds, err := a.Importer.SeriesImages(tvdbid)
			if err != nil {
				return fae.Wrap(err, "importer series images")
			}

			if len(covers) > 0 {
				eg.Go(func() error {
					err := mediumImage(series, "cover", covers[0], posterRatio)
					if err != nil {
						a.Log.Errorf("series %s cover: %v", series.ID.Hex(), err)
					}
					return nil
				})
			}
			if len(backgrounds) > 0 {
				eg.Go(func() error {
					err := mediumImage(series, "background", backgrounds[0], backgroundRatio)
					if err != nil {
						a.Log.Errorf("series %s background: %v", series.ID.Hex(), err)
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

	if err := a.DB.Medium.Save(series); err != nil {
		return fae.Wrapf(err, "saving series: %s", series.Title)
	}

	return nil
}

// type SeriesImage struct {
// 	minion.WorkerDefaults[*SeriesImage]
// 	ID    string
// 	Type  string
// 	Path  string
// 	Ratio float32
// }
//
// func (j *SeriesImage) Kind() string { return "SeriesImage" }
// func (j *SeriesImage) Work(ctx context.Context, job *minion.Job[*SeriesImage]) error {
// 	id := job.Args.ID
// 	t := job.Args.Type
// 	remote := job.Args.Path
// 	ratio := job.Args.Ratio
//
// 	series := &Series{}
// 	if err := app.DB.Series.Find(id, series); err != nil {
// 		return fae.Wrap(err, "finding series")
// 	}
//
// 	if err := seriesImage(series, t, remote, ratio); err != nil {
// 		return fae.Wrap(err, "series image")
// 	}
//
// 	if err := app.DB.Series.Save(series); err != nil {
// 		return fae.Wrap(err, "saving series")
// 	}
//
// 	return nil
// }
//
// func seriesImage(series *Series, t string, remote string, ratio float32) error {
// 	extension := filepath.Ext(remote)
// 	if len(extension) > 0 && extension[0] == '.' {
// 		extension = extension[1:]
// 	}
// 	local := fmt.Sprintf("series-%s/%s", series.ID.Hex(), t)
// 	dest := fmt.Sprintf("%s/%s.%s", app.Config.DirectoriesImages, local, extension)
// 	thumb := fmt.Sprintf("%s/%s_thumb.%s", app.Config.DirectoriesImages, local, extension)
//
// 	if err := imageDownload(remote, dest); err != nil {
// 		return fae.Wrap(err, "downloading image")
// 	}
//
// 	height := 400
// 	width := int(float32(height) * ratio)
// 	if err := imageResize(dest, thumb, width, height); err != nil {
// 		return fae.Wrap(err, "resizing image")
// 	}
//
// 	var img *Path
// 	switch t {
// 	case "cover":
// 		img = series.GetCover()
// 	case "background":
// 		img = series.GetBackground()
// 	}
//
// 	if img == nil {
// 		img = &Path{}
// 		series.Paths = append(series.Paths, img)
// 	}
//
// 	img.Type = primitive.Symbol(t)
// 	img.Remote = remote
// 	img.Local = local
// 	img.Extension = extension
//
// 	return nil
// }

func dateFromString(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Unix(0, 0)
	}
	return t
}

type EpisodeMap struct {
	byID  map[int64]*Episode
	bySE  map[int]map[int]*Episode
	byAbs map[int]*Episode
}

func (em *EpisodeMap) Add(e *Episode) error {
	sid, err := strconv.ParseInt(e.SourceID, 10, 64)
	if err != nil {
		return fae.Wrap(err, "converting source id")
	}
	em.byID[sid] = e
	if em.bySE[e.SeasonNumber] == nil {
		em.bySE[e.SeasonNumber] = map[int]*Episode{}
	}
	em.bySE[e.SeasonNumber][e.EpisodeNumber] = e
	em.byAbs[e.AbsoluteNumber] = e
	return nil
}
func (em *EpisodeMap) Find(id int64, season, episode, absolute int) *Episode {
	if id != 0 {
		if e, ok := em.byID[id]; ok {
			return e
		}
	}
	if season != 0 && episode != 0 {
		if s, ok := em.bySE[season]; ok {
			if e, ok := s[episode]; ok {
				return e
			}
		}
	}
	if absolute != 0 {
		if e, ok := em.byAbs[absolute]; ok {
			return e
		}
	}
	return nil
}

func episodeMap(id string) (*EpisodeMap, error) {
	episodeMap := &EpisodeMap{
		byID:  map[int64]*Episode{},
		bySE:  map[int]map[int]*Episode{},
		byAbs: map[int]*Episode{},
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fae.Wrap(err, "converting id")
	}

	episodes, err := app.DB.Episode.Query().Where("series_id", oid).Limit(-1).Run()
	if err != nil {
		return nil, fae.Wrap(err, "querying episodes")
	}

	for _, e := range episodes {
		episodeMap.Add(e)
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
