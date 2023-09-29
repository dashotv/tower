package app

import "github.com/dashotv/tower/internal/reader"

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
	p := reader.New(feed.Type, feed.Url)
	return p.Parse()
}
