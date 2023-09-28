package parser

import (
	"net/url"
)

func NewGeekParser(key, URL string) *GeekParser {
	return &GeekParser{
		Key: key,
		BaseParser: &BaseParser{
			URL: URL,
		},
	}
}

type GeekParser struct {
	*BaseParser
	Key string
	p   *RSSParser
}

func (p *GeekParser) Parse() error {
	u, err := url.Parse(p.URL)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Add("apikey", p.Key)
	u.RawQuery = q.Encode()

	p.p = NewRSSParser(u.String())
	return p.p.Parse()
}

func (p *GeekParser) Items() ([]Item, error) {
	return p.p.Items()
}

func (p *GeekParser) Process() error {
	return p.p.Process()
}
