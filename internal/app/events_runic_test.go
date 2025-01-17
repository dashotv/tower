package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	runic "github.com/dashotv/runic/client"
)

func TestOnRunicRelease(t *testing.T) {
	err := setupRunic(app)
	require.NoError(t, err)

	err = startWant(context.Background(), app)
	require.NoError(t, err)

	ctx := context.Background()
	req := &runic.ReleasesIndexRequest{Limit: 10}
	v, err := app.Runic.Releases.Index(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, v.Result)

	for _, r := range v.Result {
		err := onRunicReleases(app, r)
		require.NoError(t, err)
	}
}
