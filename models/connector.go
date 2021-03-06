package models

import (
	"fmt"

	"github.com/dashotv/tower/config"
)

type Connector struct {
	Download *DownloadStore
	Medium   *MediumStore
	Release  *ReleaseStore
}

var cfg *config.Config

func NewConnector() (*Connector, error) {
	cfg = config.Instance()
	var s *config.Connection
	var err error

	s, err = settingsFor("download")
	if err != nil {
		return nil, err
	}

	download, err := NewDownloadStore(s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = settingsFor("medium")
	if err != nil {
		return nil, err
	}

	medium, err := NewMediumStore(s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = settingsFor("release")
	if err != nil {
		return nil, err
	}

	release, err := NewReleaseStore(s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	c := &Connector{
		Download: download,
		Medium:   medium,
		Release:  release,
	}

	return c, nil
}

func settingsFor(name string) (*config.Connection, error) {
	if cfg.Connections["default"] == nil {
		return nil, fmt.Errorf("no connection configuration for %s", name)
	}

	if _, ok := cfg.Connections[name]; !ok {
		return cfg.Connections["default"], nil
	}

	s := cfg.Connections["default"]
	a := cfg.Connections[name]

	if a.URI != "" {
		s.URI = a.URI
	}
	if a.Database != "" {
		s.Database = a.Database
	}
	if a.Collection != "" {
		s.Collection = a.Collection
	}

	return s, nil
}
