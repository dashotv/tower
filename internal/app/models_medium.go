package app

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
)

var regexPathTv = regexp.MustCompile(`(?i)(?P<season>\d+)x(?P<episode>\d+)`)
var regexPathAnime = regexp.MustCompile(`(?i)(?P<season>\d+)x(?P<episode>\d+)(?:\s+#(?P<absolute>\d+))*`)

func (m *Medium) Destination() string {
	return filepath.Join(string(m.Kind), m.Directory)
}

// AddPathByFullpath adds a path to the medium by the full path of the file. it ensures
// that the path has a unique id and returns the path.
func (m *Medium) AddPathByFullpath(file string) *Path {
	d := filepath.Join(app.Config.DirectoriesCompleted, m.Destination())
	dest := strings.Replace(file, d+"/", "", 1)
	ext := Extension(file)

	path, ok := lo.Find(m.Paths, func(p *Path) bool {
		return p.Local == dest && p.Extension == ext
	})
	if ok && path != nil {
		if path.ID == primitive.NilObjectID {
			path.ID = primitive.NewObjectID()
		}
		return path
	}

	path = &Path{
		ID:        primitive.NewObjectID(),
		Local:     dest,
		Extension: ext,
		Type:      primitive.Symbol(fileType(fmt.Sprintf("%s.%s", dest, ext))),
	}
	m.Paths = append(m.Paths, path)
	return path
}

func (c *Connector) MediumByFile(f *File) (*Medium, error) {
	kind, name, file := f.Parts()
	filename, _ := filenameSplit(file)
	path := fmt.Sprintf("%s/%s/%s", kind, name, filename)

	if m, err := c.Medium.Query().Where("paths.local", path).Run(); err != nil {
		return nil, err
	} else if len(m) > 0 {
		return m[0], nil
	}

	switch kind {
	case "movies", "movies3d", "movies4k", "movies4h", "kids":
		return c.MediumByFilePartsMovie(kind, name)
	case "tv":
		return c.MediumByFilePartsTv(kind, name, file)
	case "anime", "donghua", "ecchi":
		return c.MediumByFilePartsAnime(kind, name, file)
	default:
		return nil, fae.Errorf("unknown kind: %s", kind)
	}
}
func (c *Connector) MediumBy(kind, name, file, ext string) (*Medium, bool, error) {
	if list, err := c.Medium.Query().Where("paths.local", fmt.Sprintf("%s/%s/%s", kind, name, file)).Where("paths.extension", ext).Run(); err != nil {
		return nil, false, err
	} else if len(list) > 0 {
		return list[0], true, nil
	}

	switch kind {
	case "movies", "movies3d", "movies4k", "movies4h", "kids":
		m, err := c.MediumByFilePartsMovie(kind, name)
		return m, false, err
	case "tv":
		m, err := c.MediumByFilePartsTv(kind, name, file)
		return m, false, err
	case "anime", "donghua", "ecchi":
		m, err := c.MediumByFilePartsAnime(kind, name, file)
		return m, false, err
	default:
		return nil, false, fae.Errorf("unknown kind: %s", kind)
	}
}

func (c *Connector) MediumByFilePartsMovie(kind, name string) (*Medium, error) {
	list, err := c.Medium.Query().Where("_type", "Movie").Where("kind", kind).Where("directory", name).Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	if len(list) > 1 {
		return nil, fae.Errorf("more than one medium found for kind: %s, name: %s", kind, name)
	}
	return list[0], nil
}

func (c *Connector) MediumByFilePartsTv(kind, name, file string) (*Medium, error) {
	series, err := c.Series.Query().Where("kind", kind).Where("directory", name).Run()
	if err != nil {
		return nil, err
	}
	if len(series) != 1 {
		return nil, fae.Errorf("series not found for kind: %s, name: %s", kind, name)
	}

	matches := regexPathTv.FindStringSubmatch(file)
	if len(matches) != 3 {
		return nil, fae.Errorf("no matches found for file: %s: %v", file, matches)
	}

	season, _ := strconv.Atoi(matches[1])
	episode, _ := strconv.Atoi(matches[2])
	list, err := c.Medium.Query().Where("_type", "Episode").Where("series_id", series[0].ID).Where("season_number", season).Where("episode_number", episode).Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fae.Errorf("not found: %s/%s/%s: %v", kind, name, file, matches)
	}
	if len(list) > 1 {
		return nil, fae.Errorf("more than one medium found for kind: %s, name: %s", kind, name)
	}
	return list[0], nil
}

func (c *Connector) MediumByFilePartsAnime(kind, name, file string) (*Medium, error) {
	series, err := c.Series.Query().Where("kind", kind).Where("directory", name).Run()
	if err != nil {
		return nil, err
	}
	if len(series) != 1 {
		return nil, fae.Errorf("series not found: %s/%s/%s", kind, name, file)
	}

	matches := regexPathAnime.FindStringSubmatch(file)
	if len(matches) != 4 {
		return nil, fae.Errorf("no matches: %s/%s/%s: %v", kind, name, file, matches)
	}

	absolute, _ := strconv.Atoi(matches[3])
	if absolute > 0 {
		list, err := c.Medium.Query().Where("_type", "Episode").Where("series_id", series[0].ID).Where("absolute_number", absolute).Run()
		if err != nil {
			return nil, err
		}
		if len(list) > 1 {
			c.Log.Warnf("more than one: %s/%s/%s: %d %d %+v", kind, name, file, absolute, list)
			return nil, fae.Errorf("more than one: %s/%s/%s: %v", kind, name, file, matches)
		}
		if len(list) == 1 {
			return list[0], nil
		}
	}

	season, _ := strconv.Atoi(matches[1])
	if season != 1 {
		return c.MediumByFilePartsTv(kind, name, file)
	}

	// absolute didn't work, try episode as absolute
	episode, _ := strconv.Atoi(matches[2])
	if episode == 0 {
		return nil, fae.Errorf("episode == 0: %s/%s/%s: %v", kind, name, file, matches)
	}

	list, err := c.Medium.Query().Where("_type", "Episode").Where("series_id", series[0].ID).Where("absolute_number", episode).Run()
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return c.MediumByFilePartsTv(kind, name, file)
	}
	if len(list) > 1 {
		c.Log.Warnf("more than one: %s/%s/%s: %d %d %+v", kind, name, file, absolute, episode, list)
		return nil, fae.Errorf("more than one: %s/%s/%s: %v", kind, name, file, matches)
	}
	return list[0], nil
}

func Background(m Medium) string {
	for _, p := range m.Paths {
		if p.Type == "background" {
			return p.Local
		}
	}
	return ""
}

func Cover(m Medium) string {
	for _, p := range m.Paths {
		if p.Type == "cover" {
			return p.Local
		}
	}
	return ""
}
