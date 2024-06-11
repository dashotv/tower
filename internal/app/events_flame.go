package app

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"

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

func onFlameCombined(app *Application, c *FlameCombined) (*EventDownloading, error) {
	list, err := app.DB.ActiveDownloads()
	if err != nil {
		return nil, fae.Wrap(err, "getting active downloads")
	}

	hashes := make(map[string]int)

	for i, d := range list {
		if len(c.Torrents) > 0 {
			t := lo.Filter(c.Torrents, func(torrent *qbt.TorrentJSON, _ int) bool {
				return strings.ToLower(torrent.Hash) == strings.ToLower(d.Thash)
			})

			if len(t) > 0 {
				d.Torrent = t[0]
				d.TorrentState = t[0].State
				if t[0].Queue > 0 {
					d.Queue = t[0].Queue
				}
				d.Progress = t[0].Progress
				if t[0].Finish > 0 && t[0].Finish != 8640000 {
					d.Eta = time.Now().Add(time.Duration(t[0].Finish) * time.Second).Format(time.RFC3339)
				}

				if d.Multi && len(d.Files) > 0 && len(t[0].Files) > 0 {
					completed := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
						tf := t[0].Files[file.Num]
						return !file.MediumID.IsZero() && tf.Progress == 100
					})
					d.FilesCompleted = len(completed)

					selected := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
						return !file.MediumID.IsZero()
					})
					d.FilesSelected = len(selected)

					wanted := lo.Filter(t[0].Files, func(file *qbt.TorrentFile, _ int) bool {
						return file.Priority > 0 && file.Progress < 100
					})
					d.FilesWanted = len(wanted)
				}
			}
		}
		if len(c.Nzbs) > 0 && d.Torrent == nil {
			n := lo.Filter(c.Nzbs, func(nzb nzbget.Group, _ int) bool {
				return fmt.Sprintf("%d", nzb.ID) == d.Thash
			})

			if len(n) > 0 {
				s := 0
				if c.NzbStatus.DownloadRate > 0 {
					s = ((n[0].RemainingSizeMB * 1024) / (c.NzbStatus.DownloadRate / 1024)) * 1000
				}
				d.Queue = float64(n[0].ID)
				d.Progress = 100.0 - (float64(n[0].RemainingSizeMB)/float64(n[0].FileSizeMB))*100.0
				if s > 0 {
					d.Eta = time.Now().Add(time.Duration(s) * time.Second).Format(time.RFC3339)
				}
			}
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

		if d.Thash != "" {
			hashes[d.Thash] = i
		}
	}

	event := &EventDownloading{
		Downloads: list,
		Hashes:    hashes,
		Metrics:   c.Metrics,
	}
	app.Cache.Set("downloads", list)
	return event, nil
}
