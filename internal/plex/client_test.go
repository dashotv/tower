package plex

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Printf("error loading .env file: %v\n", err)
	}
}

func testClient() *Client {
	return New(&ClientOptions{
		URL:              os.Getenv("PLEX_SERVER_URL"),
		Token:            os.Getenv("PLEX_TOKEN"),
		ClientIdentifier: "dashotv-test",
		Device:           "dashotv-test",
		Debug:            true,
	})
}

func TestGetAccounts(t *testing.T) {
	c := testClient()
	list, err := c.GetAccounts()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.NotEmpty(t, list)
	for _, a := range list {
		fmt.Printf("account: %+v\n", a)
	}
}

func TestGetAccount(t *testing.T) {
	c := testClient()
	account, err := c.GetAccount(2766875)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	fmt.Printf("account: %+v\n", account)
}
