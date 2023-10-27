package app

type CreateRequest struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Source string `json:"source"`
	Title  string `json:"title"`
	Date   string `json:"date"`
}
