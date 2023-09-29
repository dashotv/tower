package parser

import "github.com/dashotv/tower/internal/reader"

type Result struct {
	Raw string
}

func New(item reader.Item) *Result {
	return &Result{
		Raw: item.Title(),
	}
}
