package parser

import (
	"fmt"

	ptn "github.com/middelink/go-parse-torrent-name"
	"github.com/mmcdole/gofeed"
)

func NewRSSParser(URL string) *RSSParser {
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

func (p *RSSParser) Process() error {
	err := p.Parse()
	if err != nil {
		return err
	}

	items, err := p.Items()
	if err != nil {
		return err
	}
	for _, i := range items {
		info, err := ptn.Parse(i.Title())
		if err != nil {
			return err
		}
		fmt.Println(i.Title())
		fmt.Printf("%#v\n", info)
		fmt.Print("\n\n")
	}
	return nil
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
