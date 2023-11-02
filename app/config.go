package app

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var cfg *Config

func setupConfig() (err error) {
	cfg = &Config{}

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath("..")
	viper.AddConfigPath("../etc")
	viper.AddConfigPath(home)
	viper.AddConfigPath("/etc/tower")
	viper.SetConfigName(".tower")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// fmt.Printf("WARN: unable to read config: %s\n", err)
		return nil //errors.Wrap(err, "unable to read config")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return errors.Wrap(err, "failed to unmarshal configuration file")
	}

	if err := cfg.Validate(); err != nil {
		return errors.Wrap(err, "failed to validate config")
	}

	return nil
}

type Config struct {
	Mode        string                 `yaml:"mode"`
	Logger      string                 `yaml:"logger"`
	Port        int                    `yaml:"port"`
	Connections map[string]*Connection `yaml:"connections"`
	Cron        bool                   `yaml:"cron"`
	Auth        bool                   `yaml:"auth"`
	Plex        string                 `yaml:"plex"`
	Nats        struct {
		URL string `yaml:"url"`
	} `yaml:"nats"`
	Minion struct {
		Concurrency int `yaml:"concurrency"`
	} `yaml:"minion"`
	Redis struct {
		Address string `yaml:"address"`
	} `yaml:"redis"`
	Filesystems struct {
		Enabled     bool     `yaml:"enabled"`
		Directories []string `yaml:"directories"`
	} `yaml:"filesystems"`
}

type Connection struct {
	URI        string `yaml:"uri,omitempty"`
	Database   string `yaml:"database,omitempty"`
	Collection string `yaml:"collection,omitempty"`
}

func (c *Config) Validate() error {
	if err := c.validateDefaultConnection(); err != nil {
		return err
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
