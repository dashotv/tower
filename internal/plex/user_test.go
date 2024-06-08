package plex

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	c := testClient()
	user, err := c.GetUser(os.Getenv("PLEX_TOKEN"))
	assert.NoError(t, err)
	assert.NotNil(t, user)
	fmt.Printf("user: %+v\n", user)
}

func TestGetServicesUser(t *testing.T) {
	c := testClient()
	user, err := c.GetServicesUser(os.Getenv("PLEX_TOKEN"))
	assert.NoError(t, err)
	assert.NotNil(t, user)
	fmt.Printf("user: %+v\n", user)
}
