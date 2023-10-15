package app

import (
	"errors"
)

type Config struct {
	Mode        string                 `yaml:"mode"`
	Port        int                    `yaml:"port"`
	Connections map[string]*Connection `yaml:"connections"`
	Cron        bool                   `yaml:"cron"`
	Auth        bool                   `yaml:"auth"`
	Plex        string                 `yaml:"plex"`
	Redis       struct {
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
