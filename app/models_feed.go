package app

import "github.com/dashotv/fae"

func (c *Connector) FeedGet(id string) (*Feed, error) {
	feed, err := c.Feed.Get(id, &Feed{})
	if err != nil {
		return nil, err
	}

	// if err := c.processFeeds([]*Feed{feed}); err != nil {
	// 	return nil, err
	// }

	return feed, nil
}

func (c *Connector) FeedList(page, limit int) ([]*Feed, error) {
	skip := (page - 1) * limit
	list, err := c.Feed.Query().Desc("created_at").Limit(limit).Skip(skip).Run()
	if err != nil {
		return nil, fae.Wrap(err, "query failed")
	}

	// if err := c.processFeeds(list); err != nil {
	// 	return nil, fae.Wrap(err, "process feeds failed")
	// }

	return list, nil
}

func (c *Connector) ProcessFeeds() error {
	// feeds, err := c.Feed.Query().Where("active", true).Run()
	// if err != nil {
	// 	return err
	// }

	// for _, feed := range feeds {
	// 	if err := c.ProcessFeed(feed); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// func (c *Connector) ProcessFeed(feed *Feed) error {
// 	p := reader.New(feed.Type, feed.Url)
// 	return p.Parse()
// }

func (c *Connector) FeedUpdate(id string, data *Feed) error {
	f := &Feed{}
	err := c.Feed.Find(id, f)
	if err != nil {
		return err
	}

	f.Name = data.Name
	f.Url = data.Url
	f.Source = data.Source
	f.Type = data.Type
	f.Active = data.Active

	return c.Feed.Update(f)
}

func (c *Connector) FeedSetting(id, setting string, value bool) error {
	f := &Feed{}
	err := c.Feed.Find(id, f)
	if err != nil {
		return err
	}

	switch setting {
	case "active":
		f.Active = value
	}

	return c.Feed.Update(f)
}
