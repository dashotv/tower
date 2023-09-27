package app

import "github.com/dashotv/tower/internal/parser"

func (c *Connector) ProcessFeeds() error {
	feeds, err := c.Feed.Query().Where("active", true).Run()
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		if err := c.ProcessFeed(feed); err != nil {
			return err
		}
	}

	return nil
}

func (c *Connector) ProcessFeed(feed *Feed) error {
	p := parser.New(feed.Type, feed.Url)
	return p.Parse()
}
