package app

import (
	"os"

	"github.com/dashotv/grimoire"
)

func testConnector() *Connector {
	url := os.Getenv("TEST_MONGODB_URL")
	if url == "" {
		return nil
	}

	download, err := grimoire.New[*Download](url, "seer_development", "downloads")
	if err != nil {
		panic(err)
	}
	episode, err := grimoire.New[*Episode](url, "seer_development", "media")
	if err != nil {
		panic(err)
	}
	feed, err := grimoire.New[*Feed](url, "torch_development", "feeds")
	if err != nil {
		panic(err)
	}
	medium, err := grimoire.New[*Medium](url, "seer_development", "media")
	if err != nil {
		panic(err)
	}
	movie, err := grimoire.New[*Movie](url, "seer_development", "media")
	if err != nil {
		panic(err)
	}
	path, err := grimoire.New[*Path](url, "seer_development", "paths")
	if err != nil {
		panic(err)
	}
	release, err := grimoire.New[*Release](url, "torch_development", "torrents")
	if err != nil {
		panic(err)
	}
	series, err := grimoire.New[*Series](url, "seer_development", "media")
	if err != nil {
		panic(err)
	}
	watch, err := grimoire.New[*Watch](url, "seer_development", "watches")
	if err != nil {
		panic(err)
	}

	c := &Connector{
		Download: download,
		Episode:  episode,
		Feed:     feed,
		Medium:   medium,
		Movie:    movie,
		Path:     path,
		Release:  release,
		Series:   series,
		Watch:    watch,
	}
	return c
}
