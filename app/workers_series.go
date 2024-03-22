package app

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

type SeriesUpdate struct {
	minion.WorkerDefaults[*SeriesUpdate]
	ID         string `bson:"id" json:"id"`
	Title      string `bson:"title" json:"title"`
	SkipImages bool   `bson:"skip_images" json:"skip_images"`
}

func (j *SeriesUpdate) Kind() string { return "series_update" }
func (j *SeriesUpdate) Work(ctx context.Context, job *minion.Job[*SeriesUpdate]) error {
	id := job.Args.ID

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
		return fmt.Errorf("converting source id: %w", err)
	}

	s, err := app.Importer.Series(tvdbid)
	if err != nil {
		return fmt.Errorf("importer series: %w", err)
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

	order := importer.EpisodeOrderDefault
	anime := isAnimeKind(string(series.Kind))
	if anime {
		order = importer.EpisodeOrderAbsolute
	}

	eps, err := app.Importer.SeriesEpisodes(tvdbid, order)
	if err != nil {
		return fmt.Errorf("importer series episodes: %w", err)
	}

	episodeMap, err := episodeMap(id)
	if err != nil {
		return fmt.Errorf("building episode map: %w", err)
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
			return errors.Wrap(err, fmt.Sprintf("updating episode %s %d/%d", id, episode.SeasonNumber, episode.EpisodeNumber))
		}
	}

	all := lo.Keys(episodeMap)
	missing, updated := lo.Difference(all, found)
	if _, err := app.DB.Episode.Collection.UpdateMany(context.Background(), bson.M{"_type": "Episode", "series_id": series.ID, "source_id": bson.M{"$in": missing}}, bson.M{"$set": bson.M{"missing": time.Now()}}); err != nil {
		return fmt.Errorf("missing: %w", err)
	}
	if _, err := app.DB.Episode.Collection.UpdateMany(context.Background(), bson.M{"_type": "Episode", "series_id": series.ID, "source_id": bson.M{"$in": updated}}, bson.M{"$set": bson.M{"missing": nil}}); err != nil {
		return fmt.Errorf("found: %w", err)
	}
	if _, err := app.DB.Episode.Collection.DeleteMany(context.Background(), bson.M{"_type": "Episode", "series_id": series.ID, "missing": bson.M{"$ne": nil}, "paths.type": bson.M{"$ne": "video"}}); err != nil {
		return fmt.Errorf("missing delete: %w", err)
	}

	if err := app.DB.Series.Save(series); err != nil {
		return fmt.Errorf("saving series: %w", err)
	}

	if !job.Args.SkipImages {
		covers, backgrounds, err := app.Importer.SeriesImages(tvdbid)
		if err != nil {
			return fmt.Errorf("importer series images: %w", err)
		}

		if len(covers) > 0 {
			if err := app.Workers.Enqueue(&SeriesImage{ID: id, Type: "cover", Path: covers[0], Ratio: posterRatio}); err != nil {
				return fmt.Errorf("enqueue: cover: %w", err)
			}
		}
		if len(backgrounds) > 0 {
			if err := app.Workers.Enqueue(&SeriesImage{ID: id, Type: "background", Path: backgrounds[0], Ratio: backgroundRatio}); err != nil {
				return fmt.Errorf("enqueue: background: %w", err)
			}
		}
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
		return errors.Wrap(err, "finding series")
	}

	extension := filepath.Ext(remote)[1:]
	local := fmt.Sprintf("series-%s/%s", id, t)
	dest := fmt.Sprintf("%s/%s.%s", app.Config.DirectoriesImages, local, extension)
	thumb := fmt.Sprintf("%s/%s_thumb.%s", app.Config.DirectoriesImages, local, extension)

	if err := imageDownload(remote, dest); err != nil {
		return errors.Wrap(err, "downloading image")
	}

	height := 400
	width := int(float32(height) * ratio)
	if err := imageResize(dest, thumb, width, height); err != nil {
		return errors.Wrap(err, "resizing image")
	}

	var img *Path
	for _, p := range series.Paths {
		if string(p.Type) == t {
			img = p
			break
		}
	}

	if img == nil {
		app.Log.Info("path not found")
		img = &Path{}
		series.Paths = append(series.Paths, img)
	}

	img.Type = primitive.Symbol(t)
	img.Remote = remote
	img.Local = local
	img.Extension = extension

	if err := app.DB.Series.Update(series); err != nil {
		return errors.Wrap(err, "updating series")
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
		return nil, errors.Wrap(err, "converting id")
	}

	episodes, err := app.DB.Episode.Query().Where("series_id", oid).Limit(-1).Run()
	if err != nil {
		return nil, errors.Wrap(err, "querying episodes")
	}

	for _, e := range episodes {
		sid, err := strconv.ParseInt(e.SourceId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("converting source id: %w", err)
		}
		episodeMap[sid] = e
	}

	return episodeMap, nil
}

func episodeSEMap(id string, anime bool) (map[int]map[int]*Episode, error) {
	episodeMap := map[int]map[int]*Episode{}
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
