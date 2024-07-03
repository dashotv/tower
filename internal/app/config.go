package app

import (
	"strings"

	"github.com/caarlos0/env/v10"

	"github.com/dashotv/fae"
)

func setupConfig(app *Application) error {
	app.Config = &Config{}
	if err := env.Parse(app.Config); err != nil {
		return fae.Wrap(err, "parsing config")
	}

	if err := app.Config.Validate(); err != nil {
		return fae.Wrap(err, "failed to validate config")
	}

	return nil
}

type Config struct {
	Production            bool     `env:"PRODUCTION" envDefault:"false"`
	Logger                string   `env:"LOGGER" envDefault:"dev"`
	Port                  int      `env:"PORT" envDefault:"9000"`
	FlameURL              string   `env:"FLAME_URL"`
	ScryURL               string   `env:"SCRY_URL"`
	RunicURL              string   `env:"RUNIC_URL"`
	Plex                  string   `env:"PLEX"`
	PlexUsername          string   `env:"PLEX_USERNAME"`
	PlexToken             string   `env:"PLEX_TOKEN"`
	PlexAppName           string   `env:"PLEX_APP_NAME"`
	PlexMachineIdentifier string   `env:"PLEX_MACHINE_IDENTIFIER"`
	PlexClientIdentifier  string   `env:"PLEX_CLIENT_IDENTIFIER"`
	PlexDevice            string   `env:"PLEX_DEVICE"`
	PlexServerURL         string   `env:"PLEX_SERVER_URL"`
	PlexMetaURL           string   `env:"PLEX_META_URL"`
	PlexTvURL             string   `env:"PLEX_TV_URL"`
	DirectoriesImages     string   `env:"DIRECTORIES_IMAGES"`
	DirectoriesIncoming   string   `env:"DIRECTORIES_INCOMING"`
	DirectoriesCompleted  string   `env:"DIRECTORIES_COMPLETED"`
	DirectoriesNzbget     string   `env:"DIRECTORIES_NZBGET"`
	DirectoriesMetube     string   `env:"DIRECTORIES_METUBE"`
	FanartApiKey          string   `env:"FANART_API_KEY"`
	FanartApiURL          string   `env:"FANART_API_URL"`
	TmdbToken             string   `env:"TMDB_TOKEN"`
	TmdbImages            string   `env:"TMDB_IMAGES"`
	TvdbKey               string   `env:"TVDB_KEY"`
	DownloadsPreferred    []string `env:"DOWNLOADS_PREFERRED" envSeparator:","`
	DownloadsGroups       []string `env:"DOWNLOADS_GROUPS" envSeparator:","`
	ExtensionsVideo       []string `env:"EXTENSIONS_VIDEO" envSeparator:","`
	ExtensionsAudio       []string `env:"EXTENSIONS_AUDIO" envSeparator:","`
	ExtensionsSubtitles   []string `env:"EXTENSIONS_SUBTITLES" envSeparator:","`
	ExtensionsImages      []string `env:"EXTENSIONS_IMAGES" envSeparator:","`

	ProcessRunicEvents bool `env:"PROCESS_RUNIC_EVENTS" envDefault:"true"`
	//golem:template:app/config_partial_struct
	// DO NOT EDIT. This section is managed by github.com/dashotv/golem.
	// Models (Database)
	Connections ConnectionSet `env:"CONNECTIONS,required"`

	// Cache
	RedisAddress  string `env:"REDIS_ADDRESS,required"`
	RedisDatabase int    `env:"REDIS_DATABASE" envDefault:"0"`

	// APM
	APMServiceName string `env:"ELASTIC_APM_SERVICE_NAME,required"`
	APMServerURL   string `env:"ELASTIC_APM_SERVER_URL,required"`
	APMSecretToken string `env:"ELASTIC_APM_SECRET_TOKEN" envDefault:"0"`

	// Router Auth
	Auth           bool   `env:"AUTH" envDefault:"false"`
	ClerkSecretKey string `env:"CLERK_SECRET_KEY"`
	ClerkToken     string `env:"CLERK_TOKEN"`

	// Events
	NatsURL string `env:"NATS_URL,required"`

	// Workers
	MinionConcurrency int    `env:"MINION_CONCURRENCY" envDefault:"10"`
	MinionDebug       bool   `env:"MINION_DEBUG" envDefault:"false"`
	MinionBufferSize  int    `env:"MINION_BUFFER_SIZE" envDefault:"100"`
	MinionURI         string `env:"MINION_URI,required"`
	MinionDatabase    string `env:"MINION_DATABASE,required"`
	MinionCollection  string `env:"MINION_COLLECTION,required"`

	//golem:template:app/config_partial_struct

}

func (c *Config) Extensions() []string {
	var exts []string

	exts = append(exts, c.ExtensionsVideo...)
	exts = append(exts, c.ExtensionsAudio...)
	exts = append(exts, c.ExtensionsSubtitles...)

	return exts
}

func (c *Config) Validate() error {
	list := []func() error{
		c.validateLogger,
		//golem:template:app/config_partial_validate
		// DO NOT EDIT. This section is managed by github.com/dashotv/golem.
		c.validateDefaultConnection,

		//golem:template:app/config_partial_validate

	}

	for _, fn := range list {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) validateLogger() error {
	switch c.Logger {
	case "dev", "release":
		return nil
	default:
		return fae.New("invalid logger (must be 'dev' or 'release')")
	}
}

//golem:template:app/config_partial_connection
// DO NOT EDIT. This section is managed by github.com/dashotv/golem.

func (c *Config) validateDefaultConnection() error {
	if len(c.Connections) == 0 {
		return fae.New("you must specify a default connection")
	}

	var def *Connection
	for n, c := range c.Connections {
		if n == "default" || n == "Default" {
			def = c
			break
		}
	}

	if def == nil {
		return fae.New("no 'default' found in connections list")
	}
	if def.Database == "" {
		return fae.New("default connection must specify database")
	}
	if def.URI == "" {
		return fae.New("default connection must specify URI")
	}

	return nil
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

func (c *Config) ConnectionFor(name string) (*Connection, error) {
	def, ok := c.Connections["default"]
	if !ok {
		return nil, fae.Errorf("connection for %s: no default connection found", name)
	}

	conn, ok := c.Connections[name]
	if !ok {
		return nil, fae.Errorf("no connection named '%s'", name)
	}

	if conn.URI == "" {
		conn.URI = def.URI
	}
	if conn.Database == "" {
		conn.Database = def.Database
	}
	if conn.Collection == "" {
		conn.Collection = def.Collection
	}

	return conn, nil
}

//golem:template:app/config_partial_connection
