package app

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/pkg/errors"
)

var cfg *Config

func setupConfig() (err error) {
	cfg = &Config{}
	if err := env.Parse(cfg); err != nil {
		return errors.Wrap(err, "failed to parse environment variables")
	}

	// fmt.Println("Connections:")
	// for k, v := range cfg.Connections {
	// 	fmt.Printf("  %15s: %+v\n", k, v)
	// }

	if err := cfg.Validate(); err != nil {
		fmt.Printf("validation failed: %+v\n", cfg)
		return errors.Wrap(err, "failed to validate config")
	}

	return nil
}

type Config struct {
	Connections          ConnectionSet `env:"CONNECTIONS" envKeyValSeparator:"=" envSeparator:";"`
	Mode                 string        `env:"MODE" envDefault:"dev"`
	Logger               string        `env:"LOGGER" envDefault:"dev"`
	Port                 int           `env:"PORT" envDefault:"9000"`
	Cron                 bool          `env:"CRON" envDefault:"false"`
	Auth                 bool          `env:"AUTH" envDefault:"false"`
	Plex                 string        `env:"PLEX"`
	PlexToken            string        `env:"PLEX_TOKEN"`
	PlexAppName          string        `env:"PLEX_APP_NAME"`
	PlexClientIdentifier string        `env:"PLEX_CLIENT_IDENTIFIER"`
	PlexDevice           string        `env:"PLEX_DEVICE"`
	PlexServerURL        string        `env:"PLEX_SERVER_URL"`
	PlexMetaURL          string        `env:"PLEX_META_URL"`
	PlexTvURL            string        `env:"PLEX_TV_URL"`
	NatsURL              string        `env:"NATS_URL"`
	MinionConcurrency    int           `env:"MINION_CONCURRENCY" envDefault:"10"`
	MinionDebug          bool          `env:"MINION_DEBUG" envDefault:"false"`
	MinionBufferSize     int           `env:"MINION_BUFFER_SIZE" envDefault:"100"`
	MinionURI            string        `env:"MINION_URI"`
	MinionDatabase       string        `env:"MINION_DATABASE"`
	MinionCollection     string        `env:"MINION_COLLECTION"`
	RedisAddress         string        `env:"REDIS_ADDRESS"`
	DirectoriesImages    string        `env:"DIRECTORIES_IMAGES"`
	DirectoriesIncoming  string        `env:"DIRECTORIES_INCOMING"`
	DirectoriesCompleted string        `env:"DIRECTORIES_COMPLETED"`
	FanartApiKey         string        `env:"FANART_API_KEY"`
	FanartApiURL         string        `env:"FANART_API_URL"`
	TmdbToken            string        `env:"TMDB_TOKEN"`
	TmdbImages           string        `env:"TMDB_IMAGES"`
	TvdbKey              string        `env:"TVDB_KEY"`
	ScryURL              string        `env:"SCRY_URL"`
	DownloadsPreferred   []string      `env:"DOWNLOADS_PREFERRED" envSeparator:","`
	DownloadsGroups      []string      `env:"DOWNLOADS_GROUPS" envSeparator:","`
	ExtensionsVideo      []string      `env:"EXTENSIONS_VIDEO" envSeparator:","`
	ExtensionsAudio      []string      `env:"EXTENSIONS_AUDIO" envSeparator:","`
	ExtensionsSubtitles  []string      `env:"EXTENSIONS_SUBTITLES" envSeparator:","`
	ClerkSecretKey       string        `env:"CLERK_SECRET_KEY"`
}

func (c *Config) Extensions() []string {
	var exts []string

	exts = append(exts, c.ExtensionsVideo...)
	exts = append(exts, c.ExtensionsAudio...)
	exts = append(exts, c.ExtensionsSubtitles...)

	return exts
}

type Connection struct {
	URI        string `yaml:"uri,omitempty"`
	Database   string `yaml:"database,omitempty"`
	Collection string `yaml:"collection,omitempty"`
}

func (c *Connection) UnmarshalText(text []byte) error {
	vals := strings.Split(string(text), ",")
	c.URI = vals[0]
	c.Database = vals[1]
	c.Collection = vals[2]
	return nil
}

type ConnectionSet map[string]*Connection

func (c *ConnectionSet) UnmarshalText(text []byte) error {
	*c = make(map[string]*Connection)
	for _, conn := range strings.Split(string(text), ";") {
		kv := strings.Split(conn, "=")
		vals := strings.Split(kv[1]+",,", ",")
		(*c)[kv[0]] = &Connection{
			URI:        vals[0],
			Database:   vals[1],
			Collection: vals[2],
		}
	}
	return nil
}

func (c *Config) Validate() error {
	if err := c.validateDefaultConnection(); err != nil {
		return err
	}
	if cfg.TmdbImages == "" {
		return errors.New("tmdb.images must be set")
	}
	// TODO: add more validations?
	return nil
}

func (c *Config) validateDefaultConnection() error {
	if len(c.Connections) == 0 {
		return errors.New("you must specify a default connection")
	}

	var def *Connection
	for n, c := range c.Connections {
		if n == "default" || n == "Default" {
			def = c
			break
		}
	}

	if def == nil {
		return errors.New("no 'default' found in connections list")
	}
	if def.Database == "" {
		return errors.New("default connection must specify database")
	}
	if def.URI == "" {
		return errors.New("default connection must specify URI")
	}

	return nil
}
