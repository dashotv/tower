package app

type SearchParams struct {
	Type       string `json:"type" bson:"type"`
	Verified   bool   `json:"verified" bson:"verified"`
	Group      string `json:"group" bson:"group"`
	Author     string `json:"author" bson:"author"`
	Resolution int    `json:"resolution" bson:"resolution"`
	Source     string `json:"source" bson:"source"`
	Uncensored bool   `json:"uncensored" bson:"uncensored"`
	Bluray     bool   `json:"bluray" bson:"bluray"`
}

func NewSearchParams() *SearchParams {
	return &SearchParams{}
}
