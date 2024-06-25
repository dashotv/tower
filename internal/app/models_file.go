package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/samber/lo"

	"github.com/dashotv/fae"
)

func (f *File) Parts() (string, string, string) {
	parts := strings.Split(f.Path, string(filepath.Separator))
	if len(parts) < 6 {
		return "", "", ""
	}
	return parts[3], parts[4], parts[5]
}

func (c *Connector) FileGet(id string) (*File, error) {
	m := &File{}
	err := c.File.Find(id, m)
	if err != nil {
		return nil, err
	}

	// post process here

	return m, nil
}

func (c *Connector) FileList(page, limit int) ([]*File, int64, error) {
	skip := (page - 1) * limit

	q := c.File.Query()

	total, err := q.Count()
	if err != nil {
		return nil, 0, err
	}
	list, err := q.Desc("modified_at").Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
func (c *Connector) FileMissing(page, limit int) ([]*File, int64, error) {
	skip := (page - 1) * limit

	q := c.File.Query().Where("medium_id", primitive.NilObjectID)

	total, err := q.Count()
	if err != nil {
		return nil, 0, err
	}
	list, err := q.Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
func (c *Connector) FileCount() (int64, error) {
	return c.File.Query().Count()
}

func (c *Connector) FileByPath(path string) (*File, error) {
	list, err := c.File.Query().Where("path", path).Run()
	if err != nil {
		return nil, err
	}

	if len(list) > 1 {
		return nil, fae.Errorf("more than one file found for path: %s", path)
	}

	if len(list) == 0 {
		return &File{Path: path}, nil
	}

	return list[0], nil
}

func (c *Connector) FileFindOrCreateByPath(path string) (*File, error) {
	f, err := c.FileByPath(path)
	if err != nil {
		return nil, err
	}
	if f == nil {
		f = &File{Path: path}
	}
	return f, nil
}

func (c *Connector) DirectoryFiles(media string, page, limit int) ([]*File, int64, error) {
	skip := (page - 1) * limit

	m := &Medium{}
	if err := c.Medium.Find(media, m); err != nil {
		return nil, 0, fae.Wrap(err, "finding medium")
	}

	q := c.File.Query()
	if m.Type == "Series" {
		eids, err := c.SeriesEpisodeIDs(m.ID, skip, limit)
		if err != nil {
			return nil, 0, fae.Wrap(err, "finding series episodes")
		}
		q.In("medium_id", eids)
	} else {
		q.Where("medium_id", m.ID)
	}

	total, err := q.Count()
	if err != nil {
		return nil, 0, err
	}
	list, err := q.Desc("name").Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (c *Connector) DirectoryMedia(library string, page, limit int) ([]*Directory, int64, error) {
	skip := (page - 1) * limit

	libs, err := c.Library.Query().Where("name", library).Run()
	if err != nil {
		return nil, 0, fae.Wrap(err, "finding library")
	}
	if len(libs) != 1 {
		return nil, 0, fae.Errorf("library not found: %s", library)
	}
	lib := libs[0]

	q := c.Medium.Query().Where("kind", lib.Name)

	total, err := q.Count()
	if err != nil {
		return nil, 0, fae.Wrap(err, "counting media")
	}
	media, err := q.Asc("title").Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, 0, fae.Wrap(err, "finding media")
	}

	list := lo.Map(media, func(m *Medium, _ int) *Directory {
		count, err := c.FileCountByMedium(m)
		if err != nil {
			c.Log.Errorf("counting files for medium %s: %v", m.ID.Hex(), err)
		}
		return &Directory{Name: m.Title, Path: fmt.Sprintf("%s%c%s", lib.Name, filepath.Separator, m.ID.Hex()), Count: count}
	})

	return list, total, nil
}

func (c *Connector) FileCountByMedium(m *Medium) (int64, error) {
	if m.Type == "Series" {
		eids, err := c.SeriesEpisodeIDs(m.ID, 0, 0)
		if err != nil {
			return 0, fae.Wrap(err, "finding series episodes")
		}
		return c.File.Query().In("medium_id", eids).Count()
	}
	return c.File.Query().Where("medium_id", m.ID).Count()
}

func (c *Connector) DirectoryLibraries(page, limit int) ([]*Directory, int64, error) {
	skip := (page - 1) * limit

	total, err := c.Library.Query().Count()
	if err != nil {
		return nil, 0, fae.Wrap(err, "counting libraries")
	}

	libs, err := c.Library.Query().Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, 0, fae.Wrap(err, "finding libraries")
	}

	list := lo.Map(libs, func(l *Library, _ int) *Directory {
		return &Directory{Name: l.Name, Path: l.Name, Count: l.Count}
	})

	return list, total, nil
}
