package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Parts(t *testing.T) {
	f := &File{Path: "/mnt/media/movies3d/oz the great and powerful/oz the great and powerful.mkv"}
	kind, name, file := f.Parts()
	want := []string{"movies3d", "oz the great and powerful", "oz the great and powerful.mkv"}
	assert.Equal(t, want, []string{kind, name, file})
}
