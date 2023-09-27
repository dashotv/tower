package parser

import "github.com/mmcdole/gofeed"

func NewRSSParser(URL string) Parser {
	return &RSSParser{
		BaseParser: &BaseParser{URL: URL},
		fp:         gofeed.NewParser(),
	}
}

type RSSParser struct {
	*BaseParser
	fp   *gofeed.Parser
	feed *gofeed.Feed
}

func (p *RSSParser) Parse() error {
	feed, err := p.fp.ParseURL(p.URL)
	if err != nil {
		return err
	}

	p.feed = feed
	return nil
}

func (p *RSSParser) Items() ([]Item, error) {
	var items []Item
	for _, item := range p.feed.Items {
		items = append(items, &RSSItem{item: item})
	}
	return items, nil
}

type RSSItem struct {
	item *gofeed.Item
}

func (i *RSSItem) Title() string {
	return i.item.Title
}

func (i *RSSItem) Link() string {
	return i.item.Link
}

func (i *RSSItem) Description() string {
	return i.item.Description
}
