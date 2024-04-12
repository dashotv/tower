package app

import (
	"fmt"
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

	hashes := make(map[string]string)
	downloads := make(map[string]*Downloading)

	for _, d := range list {
		g := &Downloading{
			ID:       d.ID.Hex(),
			Infohash: d.Thash,
			Multi:    d.Multi,
		}

		if !d.MediumID.IsZero() {
			g.MediumID = d.MediumID.Hex()
		}

		g.Title = d.Title
		g.Display = d.Display
		g.Cover = d.Cover
		g.Background = d.Background

		if len(c.Torrents) > 0 {
			t := lo.Filter(c.Torrents, func(torrent *qbt.TorrentJSON, _ int) bool {
				return torrent.Hash == d.Thash
			})

			if len(t) > 0 {
				g.Torrent = t[0]
				g.TorrentState = t[0].State
				g.Queue = t[0].Queue
				g.Progress = t[0].Progress
				if t[0].Finish > 0 && t[0].Finish != 8640000 {
					g.Eta = time.Now().Add(time.Duration(t[0].Finish) * time.Second).Format(time.RFC3339)
				}

				if d.Multi && len(d.Files) > 0 && len(t[0].Files) > 0 {
					// completed := lo.Filter(t[0].Files, func(file *qbt.TorrentFile, _ int) bool {
					// 	return file.Progress == 100
					// })
					completed := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
						tf := t[0].Files[file.Num]
						return !file.MediumID.IsZero() && tf.Progress == 100
					})
					g.Files.Completed = len(completed)

					selected := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
						return !file.MediumID.IsZero()
					})
					g.Files.Selected = len(selected)
				}
			}
		}
		if len(c.Nzbs) > 0 && g.Torrent == nil {
			n := lo.Filter(c.Nzbs, func(nzb nzbget.Group, _ int) bool {
				return fmt.Sprintf("%d", nzb.ID) == d.Thash
			})

			if len(n) > 0 {
				s := 0
				if c.NzbStatus.DownloadRate > 0 {
					s = ((n[0].RemainingSizeMB * 1024) / (c.NzbStatus.DownloadRate / 1024)) * 1000
				}
				g.Queue = float64(n[0].ID)
				g.Progress = 100.0 - (float64(n[0].RemainingSizeMB)/float64(n[0].FileSizeMB))*100.0
				if s > 0 {
					g.Eta = time.Now().Add(time.Duration(s) * time.Second).Format(time.RFC3339)
				}
			}
		}

		downloads[d.ID.Hex()] = g
		if d.Thash != "" {
			hashes[d.Thash] = d.ID.Hex()
		}
	}

	event := &EventDownloading{
		Downloads: downloads,
		Hashes:    hashes,
		Metrics:   c.Metrics,
	}
	return event, nil
}
