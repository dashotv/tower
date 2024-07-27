package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dashotv/flame/qbt"
)

func TestMover_MoveDownload(t *testing.T) {
	err := setupFlame(app)
	assert.NoError(t, err)
	err = startDestination(context.Background(), app)
	assert.NoError(t, err)

	did := "668b27c9573a9d191dc0c523"
	download := &Download{}
	err = app.DB.Download.Find(did, download)
	assert.NoError(t, err)

	app.DB.processDownloads([]*Download{download})

	torrent, err := app.FlameTorrent(download.Thash)
	assert.NoError(t, err)
	assert.NotNil(t, torrent)

	mover := NewMover(app.Log.Named("TESTMOVER"), download, torrent)

	moved, err := mover.Move()
	assert.NoError(t, err)

	for _, f := range moved {
		fmt.Printf("MOVED: %s\n", f.Destination)
	}
}

func TestMover_MoveDownloadOverrides(t *testing.T) {
	err := setupFlame(app)
	assert.NoError(t, err)
	err = startDestination(context.Background(), app)
	assert.NoError(t, err)

	did := "66a43ae0f6d3142430d4f2d7"
	download := &Download{}
	err = app.DB.Download.Find(did, download)
	assert.NoError(t, err)

	app.DB.processDownloads([]*Download{download})

	mover := NewMover(app.Log.Named("TESTMOVER"), download, nil)
	moved, err := mover.Move()
	assert.NoError(t, err)

	for _, f := range moved {
		fmt.Printf("MOVED: %s\n", f.Destination)
	}
}

func TestMover_Move(t *testing.T) {
	var downloads []*Download
	var torrents []*qbt.Torrent

	err := fixture("downloads", &downloads)
	assert.NoError(t, err)
	err = fixture("torrents", &torrents)
	assert.NoError(t, err)

	mover := NewMover(app.Log.Named("TESTMOVER"), downloads[0], torrents[1])
	moved, err := mover.Move()
	assert.NoError(t, err)
	assert.Len(t, moved, 0)
}

// var moverDownload = &Download{
// 	Medium: &Medium{
// 		Type: "Series",
// 	},
// 	Files: []*DownloadFile{
// 		{
// 			Num: 1,
// 			Medium: &Medium{
// 				Type: "Episode",
// 			},
// 		},
// 		{
// 			Num: 2,
// 			Medium: &Medium{
// 				Type: "Episode",
// 			},
// 		},
// 	},
// }
//
// var moverTorrent = &qbt.TorrentJSON{
// 	Hash:         "ee1b48e6eed216440e3940f5031ce5ef2bfa0fbf",
// 	Status:       0,
// 	State:        "stalledUP",
// 	Name:         "Ze Tian Ji (Way of Choices)",
// 	Size:         1052770304,
// 	Progress:     10000,
// 	Downloaded:   1061167421,
// 	Uploaded:     0,
// 	Ratio:        0,
// 	UploadRate:   0,
// 	DownloadRate: 0,
// 	Finish:       0,
// 	Label:        "",
// 	Queue:        0,
// 	Path:         "",
// 	Files: []*qbt.TorrentFile{
// 		{
// 			ID:       0,
// 			Name:     "Ze Tian Ji (Way of Choices)/Season 1/[HaxTalks] Ze Tian Ji - Way of Choices - Ep 01 Eng Sub.mkv",
// 			Size:     171581121,
// 			Progress: 100,
// 			Priority: 0,
// 		},
// 		{
// 			ID:       1,
// 			Name:     "Ze Tian Ji (Way of Choices)/Season 1/[HaxTalks] Ze Tian Ji - Way of Choices - Ep 02 Eng Sub.mkv",
// 			Size:     350331122,
// 			Progress: 100,
// 			Priority: 0,
// 		},
// 		{
// 			ID:       2,
// 			Name:     "Ze Tian Ji (Way of Choices)/Season 1/[HaxTalks] Ze Tian Ji - Way of Choices - Ep 03 Eng Sub.mkv",
// 			Size:     332868053,
// 			Progress: 100,
// 			Priority: 1,
// 		},
// 	},
// }
