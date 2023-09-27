package parser

import "fmt"

type GeekParser struct {
	*BaseParser
}

func (p *GeekParser) Parse() error {
	return fmt.Errorf("not implemented")
}

func (p *GeekParser) Items() ([]Item, error) {
	return nil, fmt.Errorf("not implemented")
}
