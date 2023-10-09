package app

import (
	"github.com/philippgille/gokv/redis"
	"go.uber.org/zap"
)

func NewCache(log *zap.SugaredLogger, options redis.Options) (*Cache, error) {
	client, err := redis.NewClient(options)
	if err != nil {
		return nil, err
	}
	return &Cache{client: &client}, nil
}

type Cache struct {
	client *redis.Client
	log    *zap.SugaredLogger
}

func (c *Cache) Set(k string, v interface{}) error {
	return c.client.Set(k, v)
}

func (c *Cache) Get(k string, v interface{}) (bool, error) {
	return c.client.Get(k, v)
}

func (c *Cache) Delete(k string) error {
	return c.client.Delete(k)
}

func (c *Cache) Fetch(k string, v interface{}, f func() (interface{}, error)) (bool, error) {
	ok, err := c.client.Get(k, v)
	// there was an error
	if err != nil {
		return ok, err
	}
	// the item was found
	if ok {
		c.log.Infof("cache: hit: %s", k)
		return ok, nil
	}

	// get the value and set it
	v, err = f()
	if err != nil {
		return false, err
	}
	c.log.Infof("cache: miss: %s", k)
	return false, c.client.Set(k, v)
}

func (c *Cache) Close() error {
	return c.client.Close()
}
