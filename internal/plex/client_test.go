package plex

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

func testClient() *Client {
	return New(&ClientOptions{
		URL:   os.Getenv("PLEX_SERVER_URL"),
		Token: os.Getenv("PLEX_TOKEN"),
		Debug: true,
	})
}
