package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestScry_Search(t *testing.T) {
	err := setupScry(app)
	require.NoError(t, err)

	list := map[string]string{
		"Fallout S01E01":             "6566ec6dea827b91443e74ef",
		"Konosuba #24":               "6591216b65e2eca6dc1e444b",
		"Gentleman in Moscow S01e02": "660921820565844d20302252",
		"Fallout s01e04":             "65ebb4915dc3b800014c11f2",
	}
	for name, id := range list {
		t.Run(name, func(t *testing.T) {
			id, err := primitive.ObjectIDFromHex(id)
			require.NoError(t, err)

			d := &Download{Status: "searching", MediumID: id, Auto: true}
			app.DB.processDownloads([]*Download{d})

			release, err := app.ScrySearchEpisode(d.Search)
			require.NoError(t, err)
			require.NotNil(t, release)

			fmt.Printf("release: %s => %s (%d) %02dx%02d [%s]\n", name, release.Name, release.Year, release.Season, release.Episode, release.Group)
		})
	}
}
