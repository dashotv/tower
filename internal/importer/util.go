package importer

import (
	"fmt"

	"github.com/dashotv/fae"
	"github.com/dashotv/tmdb"
)

func (i *Importer) TmdbID(tvdbid int64) (int, error) {
	find, err := i.Tmdb.FindByID(fmt.Sprintf("%d", tvdbid), "tvdb_id", tmdb.String("en-US"))
	if err != nil {
		return 0, fae.Wrap(err, "tmdb id")
	}
	if find.TvResults == nil || len(find.TvResults) == 0 {
		return 0, fae.New("tmdb id: can't find id")
	}

	res := find.TvResults[0].(map[string]interface{})
	found := int(res["id"].(float64))

	return found, nil
}
