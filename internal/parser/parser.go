package parser

import (
	"fmt"
	"io"
	"net/http"
)

func New(t, URL string) Parser {
	switch t {
	case "geek":
		return &GeekParser{BaseParser: &BaseParser{URL: URL}}
	case "piratebay":
		return &PiratebayParser{BaseParser: &BaseParser{URL: URL}}
	default:
		return NewRSSParser(URL)
	}
}

type Parser interface {
	Parse() error
	Items() ([]Item, error)
}

type Item interface {
	Title() string
	Link() string
	Description() string
}

type BaseParser struct {
	URL  string
	Data string
}

func (p *BaseParser) Parse() error {
	return fmt.Errorf("not implemented")
}

func (p *BaseParser) Items() ([]Item, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *BaseParser) Read() error {
	resp, err := http.Get(p.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	p.Data = string(data)
	return nil
}
