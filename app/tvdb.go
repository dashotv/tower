package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/tvdb"
)

var tvdbClient *tvdb.Client

func setupTvdb() error {
	c, err := tvdb.Login(os.Getenv("TVDB_API_KEY"))
	if err != nil {
		return err
	}
	tvdbClient = c
	return nil
}

func TvdbUpdateSeries(payload any) error {
	id := payload.(string)

	series := &Series{}
	err := db.Series.Find(id, series)
	if err != nil {
		return err
	}

	sid, err := strconv.Atoi(series.SourceId)
	if err != nil {
		return errors.Wrap(err, "converting source id")
	}

	{
		resp, err := tvdbClient.GetSeriesTranslation(int64(sid), "eng")
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

	resp, err := tvdbClient.GetSeriesBase(int64(sid))
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

	if err := db.Series.Update(series); err != nil {
		return errors.Wrap(err, "updating series")
	}

	TvdbUpdateSeriesImages(series.ID.Hex(), int64(sid))
	workers.EnqueueWithPayload("TvdbUpdateSeriesEpisodes", series.ID.Hex())

	return nil
}

func TvdbUpdateSeriesImages(id string, sid int64) error {
	{
		r, err := tvdbClient.GetSeriesArtworks(sid, tvdb.String("eng"), tvdb.Int64(int64(2)))
		if err != nil {
			return errors.Wrap(err, "getting series artworks")
		}

		if r.Data == nil {
			return errors.New("no data")
		}

		cover := r.Data.Artworks[0]
		workers.EnqueueWithPayload("TvdbUpdateSeriesImage", &ImagePayload{id, "cover", tvdb.StringValue(cover.Image), posterRatio})
	}
	{
		r, err := tvdbClient.GetSeriesArtworks(sid, tvdb.String("eng"), tvdb.Int64(int64(3)))
		if err != nil {
			return errors.Wrap(err, "getting series artworks")
		}

		if r.Data == nil {
			return errors.New("no data")
		}
		if len(r.Data.Artworks) == 0 {
			return errors.New("no artworks")
		}

		background := r.Data.Artworks[0]
		workers.EnqueueWithPayload("TvdbUpdateSeriesImage", &ImagePayload{id, "background", tvdb.StringValue(background.Image), backgroundRatio})
	}

	return nil
}

func TvdbUpdateSeriesImage(payload any) error {
	log := log.Named("TvdbUpdateSeriesImage")
	log.Info("updating series image")

	input := payload.(*ImagePayload)
	remote := input.Path // tvdb images are full urls
	extension := filepath.Ext(input.Path)[1:]
	local := fmt.Sprintf("series-%s/%s", input.ID, input.Type)
	dest := fmt.Sprintf("%s/%s.%s", cfg.Directories.Images, local, extension)
	thumb := fmt.Sprintf("%s/%s_thumb.%s", cfg.Directories.Images, local, extension)

	if err := imageDownload(remote, dest); err != nil {
		return errors.Wrap(err, "downloading image")
	}

	height := 400
	width := int(float32(height) * input.Ratio)
	if err := imageResize(dest, thumb, width, height); err != nil {
		return errors.Wrap(err, "resizing image")
	}

	series := &Series{}
	if err := db.Series.Find(input.ID, series); err != nil {
		return errors.Wrap(err, "finding movie")
	}

	var img *Path
	for _, p := range series.Paths {
		if string(p.Type) == input.Type {
			img = p
			break
		}
	}

	if img == nil {
		log.Info("path not found")
		img = &Path{}
		series.Paths = append(series.Paths, img)
	}

	img.Type = primitive.Symbol(input.Type)
	img.Remote = remote
	img.Local = local
	img.Extension = extension

	if err := db.Series.Update(series); err != nil {
		return errors.Wrap(err, "updating series")
	}

	return nil
}

func TvdbUpdateSeriesEpisodes(payload any) error {
	log := log.Named("TvdbUpdateSeriesEpisodes")
	log.Info("updating series episodes")

	id := payload.(string)

	series := &Series{}
	err := db.Series.Find(id, series)
	if err != nil {
		return errors.Wrap(err, "finding series")
	}

	sid, err := strconv.Atoi(series.SourceId)
	if err != nil {
		return errors.Wrap(err, "converting source id")
	}

	// resp, err := tvdbClient.GetSeriesExtended(int64(sid), operations.GetSeriesExtendedMetaEpisodes.ToPointer(), tvdb.Bool(true))
	resp, err := tvdbClient.GetSeriesSeasonEpisodesTranslated(int64(sid), "eng", 0, "default")
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
	log.Infof("episode map: %d", len(episodeMap))
	log.Infof("episodes: %d", len(resp.Data.Episodes))

	for _, e := range resp.Data.Episodes {
		episode := episodeMap[tvdb.Int64Value(e.SeasonNumber)][tvdb.Int64Value(e.Number)]
		if episode == nil {
			episode = &Episode{}
		}

		log.Infof("creating/updating episode %d/%d %s", tvdb.Int64Value(e.SeasonNumber), tvdb.Int64Value(e.Number), tvdb.StringValue(e.Aired))
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

		if err := db.Episode.Save(episode); err != nil {
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

	episodes, err := db.Episode.Query().Where("_type", "Episode").Where("series_id", oid).Limit(-1).Run()
	if err != nil {
		return nil, errors.Wrap(err, "querying episodes")
	}

	log.Warnf("episodes: %d", len(episodes))

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
