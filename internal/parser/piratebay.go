package parser

import "fmt"

type PiratebayParser struct {
	*BaseParser
}

func (p *PiratebayParser) Parse() error {
	return fmt.Errorf("not implemented")
}

func (p *PiratebayParser) Items() ([]Item, error) {
	return nil, fmt.Errorf("not implemented")
}
