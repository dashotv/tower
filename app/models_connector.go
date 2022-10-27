package app

import (
	"fmt"

	"github.com/dashotv/grimoire"
)

var cfg *Config

type Connector struct {
	Download *grimoire.Store[*Download]
	Episode  *grimoire.Store[*Episode]
	Medium   *grimoire.Store[*Medium]
	Movie    *grimoire.Store[*Movie]
	Release  *grimoire.Store[*Release]
	Series   *grimoire.Store[*Series]
	Watch    *grimoire.Store[*Watch]
}

func NewConnector() (*Connector, error) {
	cfg = ConfigInstance()
	var s *Connection
	var err error

	s, err = settingsFor("download")
	if err != nil {
		return nil, err
	}

	download, err := grimoire.New[*Download](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = settingsFor("episode")
	if err != nil {
		return nil, err
	}

	episode, err := grimoire.New[*Episode](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = settingsFor("medium")
	if err != nil {
		return nil, err
	}

	medium, err := grimoire.New[*Medium](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = settingsFor("movie")
	if err != nil {
		return nil, err
	}

	movie, err := grimoire.New[*Movie](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = settingsFor("release")
	if err != nil {
		return nil, err
	}

	release, err := grimoire.New[*Release](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = settingsFor("series")
	if err != nil {
		return nil, err
	}

	series, err := grimoire.New[*Series](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	s, err = settingsFor("watch")
	if err != nil {
		return nil, err
	}

	watch, err := grimoire.New[*Watch](s.URI, s.Database, s.Collection)
	if err != nil {
		return nil, err
	}

	c := &Connector{
		Download: download,
		Episode:  episode,
		Medium:   medium,
		Movie:    movie,
		Release:  release,
		Series:   series,
		Watch:    watch,
	}

	return c, nil
}

func settingsFor(name string) (*Connection, error) {
	if cfg.Connections["default"] == nil {
		return nil, fmt.Errorf("no default config while configuring %s", name)
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
