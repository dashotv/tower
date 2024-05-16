package plex

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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

func TestGetUser(t *testing.T) {
	c := testClient()
	user, err := c.GetUser(os.Getenv("PLEX_TOKEN"))
	assert.NoError(t, err)
	assert.NotNil(t, user)
	fmt.Printf("user: %+v\n", user)
}
