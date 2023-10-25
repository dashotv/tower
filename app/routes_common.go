package app

type CreateRequest struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Title  string `json:"title"`
}
