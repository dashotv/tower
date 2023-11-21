package app

type CreateRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Date        string `json:"date"`
}
