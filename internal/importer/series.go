package importer

import (
	"github.com/dashotv/fae"
	"github.com/dashotv/tvdb"
	"github.com/dashotv/tvdb/openapi/models/operations"
)

func (i *Importer) loadSeries(id int64) (*Series, error) {
	series, err := i.loadSeriesTvdb(id)
	if err != nil {
		return nil, fae.Wrap(err, "base")
	}

	if series.Language != i.Opts.Language {
		translated, err := i.Tvdb.GetSeriesTranslation(float64(id), i.Opts.Language)
		if err != nil {
			return nil, fae.Wrap(err, "translation")
		}
		series.Title = tvdb.StringValue(translated.Data.Name)
		series.Description = tvdb.StringValue(translated.Data.Overview)
	}

	series.ID = id
	return series, nil
}

func (i *Importer) loadSeriesUpdated(since int64) ([]int64, error) {
	resp, err := i.Tvdb.Updates(since, operations.ActionUpdate.ToPointer(), tvdb.Float64(0), operations.TypeSeries.ToPointer())
	if err != nil {
		return nil, err
	}
	ints := []int64{}
	for _, s := range resp.Data {
		if s.RecordID == nil {
			continue
		}
		ints = append(ints, tvdb.Int64Value(s.RecordID))
	}

	return ints, nil
}

func (i *Importer) loadSeriesTvdb(id int64) (*Series, error) {
	resp, err := i.Tvdb.GetSeriesBase(float64(id))
	if err != nil {
		return nil, err
	}

	if resp.Data == nil {
		return nil, fae.New("no data")
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
