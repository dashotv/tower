package importer

import (
	"errors"
	"fmt"

	"github.com/dashotv/tvdb"
	"github.com/dashotv/tvdb/openapi/models/operations"
)

func (i *Importer) loadSeries(id int64) (*Series, error) {
	series, err := i.loadSeriesTvdb(id)
	if err != nil {
		return nil, fmt.Errorf("base: %w", err)
	}

	if series.Language != i.Opts.Language {
		translated, err := i.Tvdb.GetSeriesTranslation(id, i.Opts.Language)
		if err != nil {
			return nil, fmt.Errorf("translation: %w", err)
		}
		series.Title = tvdb.StringValue(translated.Data.Name)
		series.Description = tvdb.StringValue(translated.Data.Overview)
	}

	series.ID = id
	return series, nil
}

func (i *Importer) loadSeriesUpdated(since int64) ([]int64, error) {
	resp, err := i.Tvdb.Updates(since, operations.ActionUpdate.ToPointer(), tvdb.Int64(1), operations.TypeSeries.ToPointer())
	if err != nil {
		return nil, err
	}

	ints := []int64{}
	for _, s := range resp.Data {
		if s.SeriesID == nil {
			continue
		}
		ints = append(ints, tvdb.Int64Value(s.SeriesID))
	}

	return ints, nil
}

func (i *Importer) loadSeriesTvdb(id int64) (*Series, error) {
	resp, err := i.Tvdb.GetSeriesBase(id)
	if err != nil {
		return nil, err
	}

	if resp.Data == nil {
		return nil, errors.New("no data")
	}

	s := &Series{
		Title:       tvdb.StringValue(resp.Data.Name),
		Description: tvdb.StringValue(resp.Data.Overview),
		Status:      tvdb.StringValue(resp.Data.Status.Name),
		Airdate:     tvdb.StringValue(resp.Data.FirstAired),
		Language:    tvdb.StringValue(resp.Data.OriginalLanguage),
	}

	return s, nil
}
