package app

import (
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

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
	list, err := q.Limit(limit).Skip(skip).Run()
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
