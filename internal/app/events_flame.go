package app

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

		// if thashNumbersRegex.MatchString(d.Thash) && len(c.Nzbs) > 0 {
		if d.IsNzb() && len(c.Nzbs) > 0 {
			n := lo.Filter(c.Nzbs, func(nzb nzbget.Group, _ int) bool {
				return fmt.Sprintf("%d", nzb.ID) == d.Thash
			})

			if len(n) > 0 {
				if err := handleNzb(d, n[0], c.NzbStatus); err != nil {
					app.Log.Errorf("error handling nzb: %v", err)
				}
			}
		}

		if d.IsTorrent() && len(c.Torrents) > 0 {
			t := lo.Filter(c.Torrents, func(torrent *qbt.TorrentJSON, _ int) bool {
				return strings.ToLower(torrent.Hash) == strings.ToLower(d.Thash)
			})
			if len(t) > 0 {
				if err := handleTorrent(d, t[0]); err != nil {
					app.Log.Errorf("error handling torrent: %v", err)
				}
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

	event := &EventDownloading{
		Downloads: list,
		Hashes:    hashes,
		Metrics:   c.Metrics,
	}
	app.Cache.Set("downloads", list)
	return event, nil
}

func handleTorrent(d *Download, t *qbt.TorrentJSON) error {
	d.Torrent = t
	d.TorrentState = t.State
	if t.Queue > 0 {
		d.Queue = t.Queue
	}
	d.Progress = t.Progress
	if t.Finish > 0 && t.Finish != 8640000 {
		d.Eta = time.Now().Add(time.Duration(t.Finish) * time.Second).Format(time.RFC3339)
	}

	// set torrent file on download files
	for _, file := range d.Files {
		file.TorrentFile = t.Files[file.Num]
	}

	if !d.Multi || len(d.Files) == 0 || len(t.Files) == 0 {
		return nil
	}

	{
		completed := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
			tf := t.Files[file.Num]
			return !file.MediumID.IsZero() && tf.Progress == 100
		})
		d.FilesCompleted = len(completed)
	}

	{
		selected := lo.Filter(d.Files, func(file *DownloadFile, _ int) bool {
			return !file.MediumID.IsZero()
		})
		d.FilesSelected = len(selected)
	}

	{
		wanted := lo.Filter(t.Files, func(file *qbt.TorrentFile, _ int) bool {
			return file.Priority > 0 && file.Progress < 100
		})
		d.FilesWanted = len(wanted)
	}

	// sort files by torrent name
	selected := lo.Filter(d.Files, func(item *DownloadFile, index int) bool {
		return item.TorrentFile != nil && item.MediumID != primitive.NilObjectID
	})
	ignored := lo.Filter(d.Files, func(item *DownloadFile, index int) bool {
		return item.TorrentFile == nil || item.MediumID == primitive.NilObjectID
	})

	slices.SortFunc(selected, func(a, b *DownloadFile) int {
		return strings.Compare(a.TorrentFile.Name, b.TorrentFile.Name)
	})
	d.Files = append(selected, ignored...)

	return nil
}

func handleNzb(d *Download, n nzbget.Group, status nzbget.Status) error {
	s := 0
	if status.DownloadRate > 0 {
		s = ((n.RemainingSizeMB * 1024) / (status.DownloadRate / 1024)) * 1000
	}
	d.Queue = float64(n.ID)
	d.Progress = 100.0 - (float64(n.RemainingSizeMB)/float64(n.FileSizeMB))*100.0
	if s > 0 {
		d.Eta = time.Now().Add(time.Duration(s) * time.Second).Format(time.RFC3339)
	}

	return nil
}
