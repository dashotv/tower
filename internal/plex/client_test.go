package plex

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func testClient() *Client {
	return New(&ClientOptions{
		URL:   os.Getenv("PLEX_SERVER_URL"),
		Token: os.Getenv("PLEX_TOKEN"),
		Debug: false,
	})
}
