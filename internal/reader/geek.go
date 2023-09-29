package reader

import (
	"net/url"
)

func NewGeekReader(key, URL string) *GeekReader {
	return &GeekReader{
		Key: key,
		BaseReader: &BaseReader{
			URL: URL,
		},
	}
}

type GeekReader struct {
	*BaseReader
	Key string
	p   *RSSReader
}

func (p *GeekReader) Parse() error {
	u, err := url.Parse(p.URL)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Add("apikey", p.Key)
	u.RawQuery = q.Encode()

	p.p = NewRSSReader(u.String())
	return p.p.Parse()
}

func (p *GeekReader) Items() ([]Item, error) {
	return p.p.Items()
}

func (p *GeekReader) Process() error {
	return p.p.Process()
}
