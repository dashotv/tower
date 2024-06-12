package app

import (
	"slices"
	"strings"

	"github.com/dashotv/fae"
	flame "github.com/dashotv/flame/client"
	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
)

type FlameCombined struct {
	Torrents  []*qbt.TorrentJSON
	Nzbs      []nzbget.Group
	NzbStatus nzbget.Status
	Metrics   *flame.Metrics
}

// var thashNumbersRegex = regexp.MustCompile(`^\d+$`)

func onFlameCombined(app *Application, c *FlameCombined) (*EventDownloading, error) {
	list, err := app.DB.ActiveDownloads()
	if err != nil {
		return nil, fae.Wrap(err, "getting active downloads")
	}

	hashes := make(map[string]int)

	for i, d := range list {
		if d.Thash == "" || d.IsMetube() {
			continue
		}
		hashes[d.Thash] = i
		app.DB.processDownloadExtra(d, c)
	}

	slices.SortFunc(list, func(a, b *Download) int {
		d := int(a.Queue - b.Queue)
		if d != 0 {
			return d
		}

		s := a.StatusIndex() - b.StatusIndex()
		if s != 0 {
			return s
		}

		return strings.Compare(a.Title, b.Title)
	})

	event := &EventDownloading{
		Downloads: list,
		Hashes:    hashes,
		Metrics:   c.Metrics,
	}

	app.Cache.Set("downloads", list)
	return event, nil
}
