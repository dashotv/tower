package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestScry_Search(t *testing.T) {
	err := setupScry(app)
	assert.NoError(t, err)

	// Fallout S01E01
	// id, err := primitive.ObjectIDFromHex("6566ec6dea827b91443e74ef")
	// Konosuba #24
	// id, err := primitive.ObjectIDFromHex("6591216b65e2eca6dc1e444b")
	// Gentleman in Moscow S01e02
	id, err := primitive.ObjectIDFromHex("660921820565844d20302252")
	assert.NoError(t, err)

	d := &Download{Status: "searching", MediumID: id, Auto: true}
	app.DB.processDownloads([]*Download{d})

	release, err := app.ScrySearchEpisode(d.Search)
	assert.NoError(t, err)
	assert.NotNil(t, release)

	fmt.Printf("release: %+v\n", release)
}
