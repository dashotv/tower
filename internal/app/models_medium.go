package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dashotv/fae"
	"github.com/dashotv/grimoire"
	"github.com/dashotv/tower/internal/plex"
)

var regexPathTv = regexp.MustCompile(`(?i)(?P<season>\d+)x(?P<episode>\d+)`)
var regexPathAnime = regexp.MustCompile(`(?i)(?P<season>\d+)x(?P<episode>\d+)(?:\s+#(?P<absolute>\d+))*`)

// this is confusing
func (m *Medium) Destination() string {
	return filepath.Join(string(m.Kind), m.Directory)
}

func (m *Medium) BaseDir() string {
	return filepath.Join(app.Config.DirectoriesCompleted, string(m.Kind))
}

func (e *Medium) ApplyOverrides() {
	if e.Overrides == nil {
		return
	}
	a := e.Overrides.Absolute()
	if a >= 0 {
		e.HasOverrides = true
		e.AbsoluteNumber = a
	}
	s := e.Overrides.Season()
	if s >= 0 {
		e.HasOverrides = true
		e.SeasonNumber = s
	}
	ep := e.Overrides.Episode()
	if ep >= 0 {
		e.HasOverrides = true
		e.EpisodeNumber = ep
	}
}
func (m *Medium) DisplayTitle() string {
	if m.Display != "" {
		return m.Display
	}
	return m.Title
}
func (m *Medium) Year() string {
	if m.ReleaseDate.IsZero() {
		return ""
	}
	return m.ReleaseDate.Format("2006")
}

func (m *Medium) FindPathByFullPath(file string) (*Path, bool) {
	local := strings.Replace(file, app.Config.DirectoriesCompleted+"/", "", 1)
	local = strings.TrimSuffix(local, filepath.Ext(file))
	ext := Extension(file)

	return lo.Find(m.Paths, func(p *Path) bool {
		return p.Local == local && p.Extension == ext
	})
}

// AddPathByFullpath adds a path to the medium by the full path of the file. it ensures
// that the path has a unique id and returns the path.
// TODO: local should just be the filename? everything else (kind, directory, etc) is derived.
func (m *Medium) AddPathByFullpath(file string) *Path {
	local := strings.Replace(file, app.Config.DirectoriesCompleted+"/", "", 1)
	local = strings.TrimSuffix(local, filepath.Ext(file))
	ext := Extension(file)

	path, ok := lo.Find(m.Paths, func(p *Path) bool {
		return p.Local == local && p.Extension == ext
	})
	if ok && path != nil {
		if path.ID == primitive.NilObjectID {
			path.ID = primitive.NewObjectID()
		}
		path.Type = primitive.Symbol(fileType(file))
		return path
	}

	path = &Path{
		ID:        primitive.NewObjectID(),
		Local:     local,
		Extension: ext,
		Type:      primitive.Symbol(fileType(file)),
	}
	path.ParseTag()

	m.Paths = append(m.Paths, path)
	return path
}

func (m *Medium) AddPathsByMetadata(metadata *plex.KeyResponseMetadata) {
	for _, media := range metadata.Media {
		for _, part := range media.Part {
			p := m.AddPathByFullpath(part.File)
			if p != nil {
				p.Size = part.Size
				p.Resolution = metadataResolution(media.VideoResolution)
			}
		}
	}
}

func (c *Connector) MediumBySearch(title string, season, episode int) (*Medium, error) {
	// 	title = path(title)
	// 	var found *Medium
	//
	// 	// {_type:{$in:["Series","Movie"]}, $or:[{directory:"alex rider"},{search:"alex rider"}]}
	// 	q := c.Medium.Query().In("_type", []string{"Series", "Movie"})
	// 	q.Or(func(qq *grimoire.QueryBuilder[*Medium]) {
	// 		qq.Where("directory", title).Where("search", title)
	// 	})
	// 	list, err := q.Run()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if len(list) == 0 {
	// 		return nil, nil
	// 	}
	//
	// 	if found == nil {
	// 		return nil, nil
	// 	}
	//
	// 	if found.Type != "Series" {
	// 		return found, nil
	// 	}
	//
	// 	list, err := c.Medium.Query().Where("series_id", found.ID).Where("completed", false).Where("downloaded", false).Where("skipped", false).Where("season_number", season).Where("episode_number", episode).Run()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// c.Log.Debugf("MediumBySearch: %d/%d: %d", season, episode, len(list))
	// 	if len(list) == 1 {
	// 		return list[0], nil
	// 	}
	//
	// 	list, err = c.Medium.Query().Where("series_id", found.ID).Where("completed", false).Where("downloaded", false).Where("skipped", false).Where("absolute_number", episode).Run()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// c.Log.Debugf("MediumBySearch: abs %d: %d", episode, len(list))
	// 	if len(list) == 1 {
	// 		return list[0], nil
	// 	}

	return nil, nil
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

func (c *Connector) MediumByPlexMedia(media *plex.Media) (*Medium, error) {
	for _, part := range media.Part {
		kind, name, file, ext, err := pathParts(part.File)
		if err != nil {
			return nil, fae.Wrap(err, "path parts")
		}

		m, _, err := c.MediumBy(kind, name, file, ext)
		if err != nil {
			return nil, fae.Wrap(err, "medium by")
		}
		if m != nil {
			return m, nil
		}
	}

	return nil, nil
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
		list, err := c.Medium.Query().Where("_type", "Episode").Where("series_id", series[0].ID).Or(func(q *grimoire.QueryBuilder[*Medium]) {
			q.Where("absolute_number", absolute).Where("overrides.absolute_number", fmt.Sprintf("%d", absolute))
		}).Run()
		if err != nil {
			return nil, err
		}
		if len(list) > 1 {
			c.Log.Warnf("more than one episode for: %s/%s/%s: %d %+v", kind, name, file, absolute, list)
			return c.MediumByFilePartsTv(kind, name, file)
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

	list, err := c.Medium.Query().Where("_type", "Episode").Where("series_id", series[0].ID).Or(func(q *grimoire.QueryBuilder[*Medium]) {
		q.Where("absolute_number", episode).Where("overrides.absolute_number", fmt.Sprintf("%d", episode))
	}).Run()
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		c.Log.Warnf("not found episode: %s/%s/%s: %d %d %+v", kind, name, file, absolute, episode, list)
		return c.MediumByFilePartsTv(kind, name, file)
	}
	if len(list) > 1 {
		c.Log.Warnf("more than one: %s/%s/%s: %d %d %+v", kind, name, file, absolute, episode, list)
		return c.MediumByFilePartsTv(kind, name, file)
	}
	return list[0], nil
}

func (m *Medium) GetCover() *Path {
	for _, p := range m.Paths {
		if p.Type == "cover" {
			return p
		}
	}
	return nil
}

func (m *Medium) GetBackground() *Path {
	for _, p := range m.Paths {
		if p.Type == "background" {
			return p
		}
	}
	return nil
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
func (c *Connector) MediumSetting(id, setting string, value bool) error {
	m := &Medium{}
	if err := c.Medium.Find(id, m); err != nil {
		return err
	}

	c.Log.Infof("series setting: %s %t", setting, value)
	switch setting {
	case "active":
		m.Active = value
	case "favorite":
		m.Favorite = value
	case "broken":
		m.Broken = value
	case "downloaded":
		m.Downloaded = value
	case "completed":
		m.Completed = value
	}

	return c.Medium.Update(m)
}

func (c *Connector) mediumIdDeletePaths(id string) error {
	m := &Medium{}
	if err := c.Medium.Find(id, m); err != nil {
		return fae.Wrap(err, "find medium")
	}
	return c.mediumDeletePaths(m)
}

func (c *Connector) mediumDeletePaths(m *Medium) error {
	paths := m.Paths
	if m.Type == "Series" {
		err := c.Episode.Query().Where("series_id", m.ID).Batch(100, func(list []*Episode) error {
			for _, e := range list {
				paths = append(paths, e.Paths...)
			}
			return nil
		})
		if err != nil {
			return fae.Wrap(err, "listing episodes")
		}
	}

	for _, p := range paths {
		if !p.Exists() {
			continue
		}
		_ = os.Remove(p.LocalPath())
	}

	return nil
}
