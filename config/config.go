package config

import (
	"sync"
)

var once sync.Once
var instance *Config

type Config struct {
	Mode        string
	Port        int
	Connections map[string]*Connection
}

type Connection struct {
	URI        string
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

func (c *Config) Validate() error {
	// Add validations for your configuration

	return nil
}

func Instance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}
