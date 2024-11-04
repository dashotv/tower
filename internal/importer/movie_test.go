package importer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LoadMovieTmdb(t *testing.T) {
	resp, err := testImporter.loadMovieTmdb(1214701)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp)

	fmt.Printf("%+v\n", resp)
}
