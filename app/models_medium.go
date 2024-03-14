package app

import (
	"fmt"
	"regexp"
	"strings"
)

var regexPathTv = regexp.MustCompile(`(?i)(?P<season>\d+)x(?P<episode>\d+)`)
var regexPathAnime = regexp.MustCompile(`(?i)(?P<season>\d+)x(?P<episode>\d+)\s+#(?P<absolute>\d+)`)

func (c *Connector) MediumByFile(f *File) (*Medium, error) {
	kind, name, file := f.Parts()
	fext := strings.Split(file, ".")
	path := fmt.Sprintf("%s/%s/%s", kind, name, fext[0])

	if m, err := c.Medium.Query().Where("paths.local", path).Run(); err != nil {
		return nil, err
	} else if len(m) > 0 {
		return m[0], nil
	}

	switch kind {
	case "movies", "movies3d", "movies4k", "movies4h":
		return c.MediumByFilePartsMovie(kind, name)
	case "tv":
		return c.MediumByFilePartsTv(kind, name, file)
	case "anime", "donghua", "ecchi":
		return c.MediumByFilePartsAnime(kind, name, file)
	default:
		return nil, fmt.Errorf("unknown kind: %s", kind)
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
		return nil, fmt.Errorf("more than one medium found for kind: %s, name: %s", kind, name)
	}
	return list[0], nil
}

func (c *Connector) MediumByFilePartsTv(kind, name, file string) (*Medium, error) {
	series, err := c.Series.Query().Where("kind", kind).Where("directory", name).Run()
	if err != nil {
		return nil, err
	}
	if len(series) != 1 {
		return nil, fmt.Errorf("series not found for kind: %s, name: %s", kind, name)
	}

	matches := regexPathTv.FindStringSubmatch(file)
	if len(matches) != 3 {
		return nil, fmt.Errorf("no matches found for file: %s: %v", file, matches)
	}

	list, err := c.Medium.Query().Where("_type", "Episode").Where("series_id", series[0].ID).Where("season_number", matches[1]).Where("episode_number", matches[2]).Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	if len(list) > 1 {
		return nil, fmt.Errorf("more than one medium found for kind: %s, name: %s", kind, name)
	}
	return list[0], nil
}

func (c *Connector) MediumByFilePartsAnime(kind, name, file string) (*Medium, error) {
	series, err := c.Series.Query().Where("kind", kind).Where("directory", name).Run()
	if err != nil {
		return nil, err
	}
	if len(series) != 1 {
		return nil, fmt.Errorf("series not found for kind: %s, name: %s", kind, name)
	}

	matches := regexPathAnime.FindStringSubmatch(file)
	if len(matches) != 4 {
		return nil, fmt.Errorf("no matches found for file: %s: %v", file, matches)
	}

	list, err := c.Medium.Query().Where("_type", "Episode").Where("series_id", series[0].ID).Where("absolute_number", matches[3]).Run()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	if len(list) > 1 {
		return nil, fmt.Errorf("more than one medium found for kind: %s, name: %s", kind, name)
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
