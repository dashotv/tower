package app

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
	"github.com/dashotv/flame/metube"
)

var titleRegex = regexp.MustCompile(`(?i)^(?:episode|chapter)`)

func Extension(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	return ext
}

func Files(d *Download) ([]string, error) {
	out := []string{}

	if d.Thash == "" || d.IsNzb() {
		return out, nil
	}
	if d.IsMetube() {
		return FilesMetube(d)
	}

	t, err := app.FlameTorrent(d.Thash)
	if err != nil {
		return nil, err
	}

	for _, f := range t.Files {
		file := fmt.Sprintf("%s/%s", app.Config.DirectoriesIncoming, f.Name)
		if f.Progress == 100 && shouldDownloadFile(f.Name) && exists(file) {
			out = append(out, file)
		}
	}

	return out, nil
}

func FilesMetube(download *Download) ([]string, error) {
	out := []string{}

	if download.Medium.Type != "Episode" {
		return nil, fae.Errorf("unsupported medium type: %s", download.Medium.Type)
	}

	history, err := app.FlameMetubeHistory()
	if err != nil {
		return nil, fae.Wrap(err, "metube history")
	}

	done, ok := lo.Find(history.Done, func(h *metube.Download) bool {
		fmt.Printf("find: %s == %s\n", h.CustomNamePrefix, download.ID.Hex())
		return h.CustomNamePrefix == download.ID.Hex()
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

		if strings.Contains(path, download.ID.Hex()) && shouldDownloadFile(path) {
			out = append(out, path)
		}

		return nil
	})
	if err != nil {
		return nil, fae.Wrap(err, "walk")
	}

	return out, nil
}

func Destination(m *Medium) (string, error) {
	switch m.Type {
	case "Series", "Movie":
		return fmt.Sprintf("%s/%s/%s", m.Kind, m.Directory, m.Directory), nil
	case "Episode":
		return destinationEpisode(m)
	default:
		return "", fae.Errorf("unknown type: %s", m.Type)
	}
}

func destinationEpisode(m *Medium) (string, error) {
	s := &Series{}
	err := app.DB.Series.FindByID(m.SeriesID, s)
	if err != nil {
		return "", err
	}

	e := &Episode{}
	err = app.DB.Episode.FindByID(m.ID, e)
	if err != nil {
		return "", err
	}

	out := ""
	if isAnimeKind(string(s.Kind)) && m.AbsoluteNumber > 0 {
		out = fmt.Sprintf("%s/%s/%s - %02dx%02d #%03d", s.Kind, s.Directory, s.Directory, m.SeasonNumber, m.EpisodeNumber, m.AbsoluteNumber)
	} else {
		out = fmt.Sprintf("%s/%s/%s - %02dx%02d", s.Kind, s.Directory, s.Directory, m.SeasonNumber, m.EpisodeNumber)
	}

	if e.Title != "" && !titleRegex.MatchString(e.Title) {
		out = fmt.Sprintf("%s - %s", out, path(e.Title))
	}

	return out, nil
}
