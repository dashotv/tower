package reader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func New(t, URL string) Reader {
	switch t {
	case "geek":
		return NewGeekReader(os.Getenv("NZBGEEK_API_KEY"), URL)
	default:
		return NewRSSReader(URL)
	}
}

type Reader interface {
	Parse() error
	Items() ([]Item, error)
	Process() error
}

type Item interface {
	Title() string
	Link() string
	Description() string
}

type BaseReader struct {
	URL  string
	Data string
}

func (p *BaseReader) Parse() error {
	return fmt.Errorf("not implemented")
}

func (p *BaseReader) Items() ([]Item, error) {
	return nil, fmt.Errorf("not implemented")
}

func (p *BaseReader) Read() error {
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
