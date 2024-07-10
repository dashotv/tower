package app

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
	"github.com/dashotv/grimoire"
	runic "github.com/dashotv/runic/client"
	"github.com/dashotv/runic/parser"
)

func init() {
	initializers = append(initializers, setupRunic)
}

func setupRunic(app *Application) error {
	app.Runic = runic.New(app.Config.RunicURL)
	return nil
}

func (a *Application) RunicFindEpisode(seriesID primitive.ObjectID, title, type_ string) (*Episode, error) {
	info, err := parser.Parse(title, type_)
	if err != nil {
		return nil, fae.Wrap(err, "parsing title")
	}
	if info.Season == 0 && info.Episode == 0 {
		return nil, nil
	}

	q := app.DB.Episode.Query().Where("series_id", seriesID).Asc("season_number").Asc("episode_number").Asc("absolute_number")
	if type_ == "anime" {
		if info.Season == 0 {
			// if season is 0, only check absolute number
			q = q.Where("absolute_number", info.Episode)
		} else {
			// if season > 0, check both absolute and season/episode number
			q = q.ComplexOr(func(qq, qr *grimoire.QueryBuilder[*Episode]) {
				qq.Where("season_number", info.Season).Where("episode_number", info.Episode)
				qr.Where("absolute_number", info.Episode)
			})
		}
	} else {
		q = q.Where("season_number", info.Season).Where("episode_number", info.Episode)
	}
	eps, err := q.Run()
	if err != nil {
		return nil, fae.Wrap(err, "querying episode")
	}
	if len(eps) == 0 {
		return nil, nil
	}
	// if len(eps) > 1 {
	// 	return nil, fae.Errorf("querying episode: multiple episodes found")
	// }

	return eps[0], nil
}
