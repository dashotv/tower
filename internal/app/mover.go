package app

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/dashotv/fae"
	"github.com/dashotv/flame/metube"
	"github.com/dashotv/flame/qbt"
)

func NewMover(log *zap.SugaredLogger, download *Download, torrent *qbt.Torrent) *Mover {
	m := &Mover{
		Log:      log,
		Download: download,
		Torrent:  torrent,
		moved:    []*MoverFile{},
	}

	m.movefunc = m.FileLink
	if !app.Config.Production {
		m.movefunc = m.testFileLink
	}

	return m
}

type Mover struct {
	Log      *zap.SugaredLogger
	Download *Download
	Torrent  *qbt.Torrent
	moved    []*MoverFile
	movefunc func(string, string, bool) error
}
type MoverFile struct {
	Source      string
	Destination string
	Medium      *Medium
}

func (m *Mover) List() ([]string, error) {
	return m.torrentList()
}

func (m *Mover) torrentList() ([]string, error) {
	out := []string{}

	for _, f := range m.Torrent.Files {
		file := fmt.Sprintf("%s/%s", app.Config.DirectoriesIncoming, f.Name)
		df := m.Download.File(f.ID)
		if f.Progress == 100 && shouldDownloadFile(f.Name) && exists(file) && df != nil && df.Medium != nil {
			out = append(out, f.Name)
		}
	}

	return out, nil
}

func (m *Mover) metubeList() ([]string, error) {
	out := []string{}

	if m.Download.Medium.Type != "Episode" {
		return nil, fae.Errorf("unsupported medium type: %s", m.Download.Medium.Type)
	}

	history, err := app.FlameMetubeHistory()
	if err != nil {
		return nil, fae.Wrap(err, "metube history")
	}

	done, ok := lo.Find(history.Done, func(h *metube.Download) bool {
		return h.CustomNamePrefix == m.Download.ID.Hex()
	})
	if !ok || done == nil {
		return nil, nil
	}

	err = filepath.WalkDir(app.Config.DirectoriesMetube, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.Contains(path, m.Download.ID.Hex()) && shouldDownloadFile(path) {
			out = append(out, path)
		}

		return nil
	})
	if err != nil {
		return nil, fae.Wrap(err, "walk")
	}

	return out, nil
}

func (m *Mover) Move() ([]*MoverFile, error) {
	if m.Download.IsMetube() {
		return m.moveMetube()
	}
	if m.Download.Medium.Type == "Series" {
		return m.moveSeries()
	}
	return m.moveFiles()
}

func (m *Mover) moveSeries() ([]*MoverFile, error) {
	dfiles := m.Download.Files
	numToDf := map[int]*DownloadFile{}
	for _, df := range dfiles {
		numToDf[df.Num] = df
	}

	tfiles := lo.Filter(m.Torrent.Files, func(f *qbt.TorrentFile, _ int) bool {
		return f.Progress == 100 && shouldDownloadFile(f.Name) && numToDf[f.ID].Medium != nil
	})

	for _, tf := range tfiles {
		medium := numToDf[tf.ID].Medium

		err := m.moveFile(fmt.Sprintf("%s/%s", app.Config.DirectoriesIncoming, tf.Name), medium)
		if err != nil {
			return nil, err
		}
	}

	return m.moved, nil
}

func (m *Mover) moveFiles() ([]*MoverFile, error) {
	files, err := m.List()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		err := m.moveFile(fmt.Sprintf("%s/%s", app.Config.DirectoriesIncoming, file), m.Download.Medium)
		if err != nil {
			return nil, err
		}
	}

	return m.moved, nil
}

func (m *Mover) moveMetube() ([]*MoverFile, error) {
	files, err := m.metubeList()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		err := m.moveFile(file, m.Download.Medium)
		if err != nil {
			return nil, err
		}
	}

	return m.moved, nil
}

func (m *Mover) moveFile(source string, medium *Medium) error {
	ext := Extension(source)

	dest, err := app.Destinator.Destination(m.Download.Kind, medium)
	if err != nil {
		return fae.Wrap(err, "getting destination")
	}

	tag := ""
	if m.Download.Tag != "" {
		tag = fmt.Sprintf(" [%s]", m.Download.Tag)
	}
	destination := fmt.Sprintf("%s%s.%s", dest, tag, ext)

	if medium == nil {
		m.Log.Debugf("skipping %s", source)
		return nil
	}

	if !exists(source) {
		return fae.Errorf("source does not exist: %s", source)
	}
	if exists(destination) {
		if !m.Download.Force {
			// notifier.Log.Warn("DownloadMove", fmt.Sprintf("destination exists, force false: %s", destination))
			return nil
		}

		match, err := sumFiles(source, destination)
		if err != nil {
			return fae.Errorf("failed to sum files")
		}
		if match {
			// notifier.Log.Warn("DownloadMove", fmt.Sprintf("destination exists, sums match: %s", destination))
			return fae.Errorf("destination exists, sums match: %s", destination)
		}
	}

	if err := m.movefunc(source, destination, m.Download.Force); err != nil {
		return fae.Wrap(err, "link")
	}

	m.moved = append(m.moved, &MoverFile{Source: source, Destination: destination, Medium: medium})
	return nil
}

func (m *Mover) FileLink(source, destination string, force bool) error {
	return FileLink(source, destination, force)
}

func (m *Mover) testFileLink(source, destination string, force bool) error {
	m.Log.Warnf("linking[%t] %s -> %s\n", force, source, destination)
	return nil
}
