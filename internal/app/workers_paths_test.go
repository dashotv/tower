package app

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestApp_PathDest(t *testing.T) {
	err := startDestination(context.Background(), app)
	require.NoError(t, err)

	m := &Medium{}
	p := "/mnt/media/donghua/stellar transformation/stellar transformation - 01x083 - 03 [animexin 1080].mp4"
	kind := primitive.Symbol("donghua")
	err = app.DB.Medium.Find("679a93502bcb0e9a4864ffef", m)
	require.NoError(t, err)
	require.NotNil(t, m)

	path, ok := lo.Find(m.Paths, func(item *Path) bool {
		return item.LocalPath() == p
	})
	require.NotNil(t, path)
	require.True(t, ok)

	err = app.pathDest(m, kind, path)
	require.NoError(t, err)
	require.False(t, path.Rename)
}
