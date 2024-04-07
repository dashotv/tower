package app

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/dashotv/grimoire"
)

func init() {
	err := appSetup()
	if err != nil {
		panic(err)
	}
}

var envVars = []string{"CONNECTIONS", "NATS_URL", "REDIS_ADDRESS", "MINION_URI", "FLAME_URL"}

func appSetup() error {
	if app != nil {
		fmt.Println("app already setup")
		return nil
	}

	err := envReplaceAll("host.docker.internal", "localhost", envVars)
	if err != nil {
		return err
	}

	app = &Application{}
	list := []func(*Application) error{
		setupConfig,
		setupLogger,
		setupEvents,
		setupDb,
	}

	for _, f := range list {
		if err := f(app); err != nil {
			return err
		}
	}

	return nil
}

func envReplaceAll(old, new string, vars []string) error {
	for _, v := range vars {
		if err := os.Setenv(v, strings.ReplaceAll(os.Getenv(v), old, new)); err != nil {
			return err
		}
	}
	return nil
}

func testConnector() *Connector {
	app := &Application{}
	setupConfig(app)
	setupLogger(app)

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
		Log:      app.Log.Named("db"),
		Download: download,
		Episode:  episode,
		Feed:     feed,
		Medium:   medium,
		Movie:    movie,
		Release:  release,
		Series:   series,
		Watch:    watch,
	}

	app.DB = c

	return c
}
