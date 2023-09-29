package parser

type Result struct {
	Raw string
}

func New(title string) *Result {
	return &Result{
		Raw: title,
	}
}
