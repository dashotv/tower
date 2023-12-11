package app

import (
	"fmt"
	"time"

	"github.com/samber/lo"

	flame "github.com/dashotv/flame/app"
	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/flame/qbt"
)

type FlameCombined struct {
	Torrents  []*qbt.TorrentJSON
	Nzbs      []nzbget.Group
	NzbStatus nzbget.Status
	Metrics   *flame.Metrics
}

type EventTowerDownloading struct {
	Downloads map[string]*Downloading `json:"downloads,omitempty"`
	Hashes    map[string]string       `json:"hashes,omitempty"`
	Metrics   *flame.Metrics          `json:"metrics,omitempty"`
}

type Downloading struct {
	ID           string           `json:"id,omitempty"`
	MediumID     string           `json:"medium_id,omitempty"`
	Multi        bool             `json:"multi,omitempty"`
	Infohash     string           `json:"infohash,omitempty"`
	Torrent      *qbt.TorrentJSON `json:"torrent,omitempty"`
	Queue        float64          `json:"queue"`
	Progress     float64          `json:"progress"`
	Eta          string           `json:"eta,omitempty"`
	TorrentState string           `json:"torrent_state,omitempty"`
	Files        struct {
		Completed int `json:"completed,omitempty"`
		Selected  int `json:"selected,omitempty"`
	} `json:"files,omitempty"`
}

func sendDownloading(c *FlameCombined) {
	list, err := db.ActiveDownloads()
	if err != nil {
		events.Log.Errorf("error getting active downloads: %s", err)
		return
	}

	hashes := make(map[string]string)
	downloads := make(map[string]*Downloading)

	for _, d := range list {
		g := &Downloading{
			ID:       d.ID.Hex(),
			Infohash: d.Thash,
			Multi:    d.Multi,
		}

		if !d.MediumId.IsZero() {
			g.MediumID = d.MediumId.Hex()
		}

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
						return !file.MediumId.IsZero() && tf.Progress == 100
					})
					g.Files.Completed = len(completed)

					selected := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
						return !file.MediumId.IsZero()
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
	event := &EventTowerDownloading{
		Downloads: downloads,
		Hashes:    hashes,
		Metrics:   c.Metrics,
	}

	// for _, d := range downloads {
	// 	log.Debugf("downloading: %6.2f %d/%d", d.Progress, d.Files.Completed, d.Files.Selected)
	// }
	events.Send("tower.downloading", event)
}
