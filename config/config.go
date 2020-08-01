package config

import (
	"sync"
)

var once sync.Once
var instance *Config

type Config struct {
	Mode string
	Port int
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
