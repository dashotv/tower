package reader

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func NewRSSReader(URL string) *RSSReader {
	return &RSSReader{
		BaseReader: &BaseReader{URL: URL},
		fp:         gofeed.NewParser(),
	}
}

type RSSReader struct {
	*BaseReader
	fp   *gofeed.Parser
	feed *gofeed.Feed
}

func (p *RSSReader) Parse() error {
	feed, err := p.fp.ParseURL(p.URL)
	if err != nil {
		return err
	}

	p.feed = feed
	return nil
}

func (p *RSSReader) Items() ([]Item, error) {
	var items []Item
	for _, item := range p.feed.Items {
		items = append(items, &RSSItem{item: item})
	}
	return items, nil
}

func (p *RSSReader) Process() error {
	err := p.Parse()
	if err != nil {
		return err
	}

	items, err := p.Items()
	if err != nil {
		return err
	}
	for _, i := range items {
		// info, err := ptn.Parse(i.Title())
		// if err != nil {
		// 	return err
		// }
		fmt.Println(i.Title())
		fmt.Print("\n\n")
	}
	return nil
}

type RSSItem struct {
	item *gofeed.Item
}

func (i *RSSItem) Title() string {
	return fmt.Sprintf("%s\n%#v\n", i.item.Title, i.item)
}

func (i *RSSItem) Link() string {
	return i.item.Link
}

func (i *RSSItem) Description() string {
	return i.item.Description
}

func (i *RSSItem) Guid() string {
	return i.item.GUID
}

func (i *RSSItem) Published() string {
	return i.item.Published
}

func (i *RSSItem) Updated() string {
	return i.item.Updated
}

//
// func (i *RSSItem) Enclosure() *Enclosure {
// 	if i.item.Enclosures == nil {
// 		return nil
// 	}
// 	return &Enclosure{
// 		URL:  i.item.Enclosures[0].URL,
// 		Type: i.item.Enclosures[0].Type,
// 	}
// }

func (i *RSSItem) Author() string {
	return i.item.Author.Name
}
