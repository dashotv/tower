package app

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
	"github.com/dashotv/grimoire"
	runic "github.com/dashotv/runic/client"
)

func init() {
	initializers = append(initializers, setupRunic)
}

func setupRunic(app *Application) error {
	app.Runic = runic.New(app.Config.RunicURL)
	return nil
}

func (a *Application) RunicFindEpisode(seriesID primitive.ObjectID, title, type_ string) (*Episode, error) {
	req := &runic.ParserTitleRequest{Title: title, Type: type_}
	resp, err := app.Runic.Parser.Title(context.Background(), req)
	if err != nil {
		return nil, fae.Wrap(err, "parsing title")
	}
	if resp == nil || resp.Result == nil {
		return nil, fae.Wrap(err, "parsing title, response nil")
	}
	if resp.Error {
		return nil, fae.Errorf("parsing title: %s", resp.Message)
	}

	info := resp.Result
	q := app.DB.Episode.Query().Where("series_id", seriesID).Asc("season_number").Asc("episode_number").Asc("absolute_number")
	if type_ == "anime" {
		q = q.Or(func(q *grimoire.QueryBuilder[*Episode]) {
			q.Where("episode_number", info.Episode)
			q.Where("absolute_number", info.Episode)
		})
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
